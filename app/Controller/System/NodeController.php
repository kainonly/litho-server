<?php
declare(strict_types=1);

namespace App\Controller\System;

use App\GrpcClient\SSHServiceInterface;
use Hyperf\Curd\Common\AddAfterParams;
use Hyperf\Curd\Common\DeleteAfterParams;
use Hyperf\Curd\Common\EditAfterParams;
use Hyperf\DbConnection\Db;
use Hyperf\Di\Annotation\Inject;
use Hyperf\Extra\Cipher\CipherInterface;
use Hyperf\Utils\Context;
use SSHMicroservice\Response;
use stdClass;

class NodeController extends BaseController
{
    /**
     * @Inject()
     * @var CipherInterface
     */
    private CipherInterface $cipher;
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
            ->originListsModel('node')
            ->setOrder('create_time', 'desc')
            ->setField(['id', 'group', 'name', 'host', 'port', 'username', 'status'])
            ->result();
    }

    public function lists(): array
    {
        $validate = $this->curd->listsValidation();
        if ($validate->fails()) {
            return [
                'error' => 1,
                'msg' => $validate->errors()
            ];
        }

        return $this->curd
            ->listsModel('node')
            ->setOrder('create_time', 'desc')
            ->setField(['id', 'group', 'name', 'host', 'port', 'username', 'status'])
            ->result();
    }

    public function get(): array
    {
        $validate = $this->curd->getValidation();
        if ($validate->fails()) {
            return [
                'error' => 1,
                'msg' => $validate->errors()
            ];
        }
        return $this->curd
            ->getModel('node')
            ->setField(['id', 'group', 'name', 'host', 'port', 'username', 'status'])
            ->result();
    }

    public function add(): array
    {
        $body = $this->request->post();
        $validate = $this->curd->addValidation([
            'name' => 'required',
            'group' => 'required',
            'host' => 'required|ip',
            'port' => 'required|integer',
            'username' => 'required|string',
            'password' => 'required_without:private_key',
            'private_key' => 'required_without:password'
        ]);
        if ($validate->fails()) {
            return [
                'error' => 1,
                'msg' => $validate->errors()
            ];
        }
        $attr = $this->preAttribute($body);
        return $this->curd
            ->addModel('node', $body)
            ->afterHook(function (AddAfterParams $params) use ($body, $attr) {
                $response = $this->afterRequest($params->getId(), $body, $attr);
                $this->sshService->close();
                if ($response->getError() !== 0) {
                    Context::set('error', [
                        'error' => 1,
                        'msg' => $response->getMsg()
                    ]);
                    return false;
                }
                return true;
            })
            ->result();
    }

    public function edit(): array
    {
        $body = $this->request->post();
        $validate = $this->curd->editValidation([
            'name' => 'required',
            'group' => 'required',
            'host' => 'required|ip',
            'port' => 'required|integer',
            'username' => 'required|string',
        ]);
        if ($validate->fails()) {
            return [
                'error' => 1,
                'msg' => $validate->errors()
            ];
        }
        $attr = $this->preAttribute($body);
        return $this->curd
            ->editModel('node', $body)
            ->afterHook(function (EditAfterParams $params) use ($body, $attr) {
                $response = $this->afterRequest($params->getId(), $body, $attr);
                $this->sshService->close();
                if ($response->getError() !== 0) {
                    Context::set('error', [
                        'error' => 1,
                        'msg' => $response->getMsg()
                    ]);
                    return false;
                }
                return true;
            })
            ->result();
    }

    public function delete(): array
    {
        $validate = $this->curd->deleteValidation();
        if ($validate->fails()) {
            return [
                'error' => 1,
                'msg' => $validate->errors()
            ];
        }
        return $this->curd
            ->deleteModel('node')
            ->afterHook(function (DeleteAfterParams $params) {
                $response = $this->sshService->delete((string)$params->getId()[0]);
                $this->sshService->close();
                if ($response->getError() !== 0) {
                    Context::set('error', [
                        'error' => 1,
                        'msg' => $response->getMsg()
                    ]);
                    return false;
                }
                return true;
            })
            ->result();
    }

    /**
     * @param array $body
     * @return stdClass
     */
    private function preAttribute(array &$body): stdClass
    {
        $data = new stdClass();
        if (!empty($body['password'])) {
            $data->password = $body['password'];
            $body['password'] = $this->cipher->encrypt($data->password);
        } else {
            unset($body['password']);
        }
        if (!empty($body['private_key'])) {
            $data->privateKey = base64_encode($body['private_key']);
            $body['private_key'] = $this->cipher->encrypt($data->privateKey);
        } else {
            unset($body['private_key']);
        }
        if (!empty($body['passphrase'])) {
            $data->passphrase = $body['passphrase'];
            $body['passphrase'] = $this->cipher->encrypt($data->passphrase);
        } else {
            unset($body['passphrase']);
        }
        return $data;
    }

    /**
     * @param int $id
     * @param array $body
     * @param stdClass $attribute
     * @return Response
     */
    private function afterRequest(int $id, array $body, stdClass $attribute): Response
    {
        return $this->sshService->put(
            (string)$id,
            $body['host'],
            $body['port'],
            $body['username'],
            $attribute->password ?? '',
            $attribute->privateKey ?? '',
            $attribute->passphrase ?? ''
        );
    }

    /**
     * 连接测试接口
     * @return array
     */
    public function connected(): array
    {
        $body = $this->request->post();
        $validator = $this->validation->make($body, [
            'host' => 'required|ip',
            'port' => 'required|integer',
            'username' => 'required|string',
            'password' => 'required_without:private_key',
            'private_key' => 'required_without:password',
        ]);

        if ($validator->fails()) {
            return [
                'error' => 1,
                'msg' => $validator->errors()
            ];
        }

        $response = $this->sshService->testing(
            $body['host'],
            $body['port'],
            $body['username'],
            $this->post['password'] ?? '',
            !empty($body['private_key']) ? base64_encode($body['private_key']) : '',
            $this->post['passphrase'] ?? ''
        );
        $this->sshService->close();
        return [
            'error' => $response->getError(),
            'msg' => $response->getMsg()
        ];
    }

    /**
     * 同步服务
     * @return array
     */
    public function sync(): array
    {
        $response = $this->sshService->all();
        if ($response->getError() !== 0) {
            return [
                'error' => 1,
                'msg' => $response->getMsg()
            ];
        }
        foreach ($response->getData()->getIterator() as $identity) {
            $deleteResponse = $this->sshService->delete((string)$identity);
            if ($deleteResponse->getError() !== 0) {
                return [
                    'error' => 1,
                    'msg' => "<$identity> " . $deleteResponse->getMsg()
                ];
            }
        }
        $query = Db::table('node')->get();
        foreach ($query->toArray() as $value) {
            $putResponse = $this->sshService->put(
                (string)$value->id,
                $value->host,
                $value->port,
                $value->username,
                !empty($value->password) ? $this->cipher->decrypt($value->password) : '',
                !empty($value->private_key) ? $this->cipher->decrypt($value->private_key) : '',
                !empty($value->passphrase) ? $this->cipher->decrypt($value->passphrase) : ''
            );
            if ($putResponse->getError() !== 0) {
                return [
                    'error' => 1,
                    'msg' => "<{$value->id}> " . $putResponse->getMsg()
                ];
            }
            $result = TunnelsController::sync($value->id, $this->sshService);
            if ($result['error'] === 1) {
                return $result;
            }
        }
        $this->sshService->close();
        return [
            'error' => 0,
            'msg' => 'ok'
        ];
    }

    /**
     * 系统负载查询
     * @return array
     */
    public function uptime(): array
    {
        $body = $this->request->post();
        $validator = $this->validation->make($body, [
            'identity' => 'required',
        ]);

        if ($validator->fails()) {
            return [
                'error' => 1,
                'msg' => $validator->errors()
            ];
        }

        $response = $this->sshService->exec((string)$body['identity'], 'uptime');
        $this->sshService->close();
        return $response->getError() !== 0 ? [
            'error' => 1,
            'msg' => $response->getMsg()
        ] : [
            'error' => 0,
            'data' => $response->getData()
        ];
    }
}