<?php
declare(strict_types=1);

namespace App\Controller\System;

use App\GrpcClient\SSHServiceInterface;
use Hyperf\DbConnection\Db;
use Hyperf\Di\Annotation\Inject;
use Hyperf\Utils\Context;

class TunnelsController extends BaseController
{
    /**
     * @Inject()
     * @var SSHServiceInterface
     */
    private SSHServiceInterface $sshService;

    public function originLists(): array
    {
        $validate = $this->curd->originListsValidation();
        if ($validate->fails()) {
            return [
                'error' => 1,
                'msg' => $validate->errors()
            ];
        }

        return $this->curd
            ->originListsModel('tunnels')
            ->result();
    }

    public function add(): array
    {
        $body = $this->request->post();
        $validate = $this->curd->addValidation([
            'node' => 'required',
            'src_ip' => 'required',
            'src_port' => 'required',
            'dst_port' => 'required'
        ]);
        if ($validate->fails()) {
            return [
                'error' => 1,
                'msg' => $validate->errors()
            ];
        }
        return $this->curd
            ->addModel('tunnels')
            ->setAutoTimestamp(false)
            ->afterHook(function () use ($body) {
                $response = self::sync($body['node'], $this->sshService);
                if ($response['error'] === 0) {
                    Context::set('error', $response);
                    return false;
                }
                return true;
            })
            ->result();
    }

    public function delete(): array
    {
        $body = $this->request->post();
        $validate = $this->curd->deleteValidation();
        if ($validate->fails()) {
            return [
                'error' => 1,
                'msg' => $validate->errors()
            ];
        }
        $data = Db::table('tunnels')
            ->where('id', '=', $body['id'][0])
            ->first();
        if (empty($data)) {
            return [
                'error' => 1,
                'msg' => 'tunnels not exists'
            ];
        }
        return $this->curd
            ->deleteModel('tunnels')
            ->afterHook(function () use ($data) {
                $response = self::sync($data->node, $this->sshService);
                if ($response['error'] === 0) {
                    Context::set('error', $response);
                    return false;
                }
                return true;
            })
            ->result();
    }

    /**
     * 同步隧道服务
     * @param int $identity
     * @return array
     */
    static function sync(int $identity, SSHServiceInterface $SSHService): array
    {
        $query = Db::table('tunnels')
            ->where('node', '=', $identity)
            ->get();

        if ($query->isEmpty()) {
            return [
                'error' => 0,
                'msg' => 'empty sync'
            ];
        }

        $tunnels = [];
        foreach ($query->toArray() as $value) {
            $tunnels[] = [
                'src_ip' => $value->src_ip,
                'src_port' => $value->src_port,
                'dst_ip' => '0.0.0.0',
                'dst_port' => $value->dst_port
            ];
        }
        $response = $SSHService->tunnels((string)$identity, $tunnels);
        return [
            'error' => $response->getError(),
            'msg' => $response->getMsg()
        ];
    }
}