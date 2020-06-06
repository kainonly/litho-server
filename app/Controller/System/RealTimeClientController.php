<?php
declare(strict_types=1);

namespace App\Controller\System;

use App\Event\RealTimeRefresh;
use Hyperf\DbConnection\Db;
use Hyperf\Di\Annotation\Inject;
use Psr\EventDispatcher\EventDispatcherInterface;

class RealTimeClientController extends BaseController
{
    /**
     * @Inject()
     * @var EventDispatcherInterface
     */
    private EventDispatcherInterface $eventDispatcher;

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
            ->originListsModel('real_time_client')
            ->setOrder('create_time', 'desc')
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
            ->listsModel('real_time_client')
            ->setOrder('create_time', 'desc')
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
            ->getModel('real_time_client')
            ->result();
    }

    public function add(): array
    {
        $validate = $this->curd->addValidation([
            'group' => 'required',
            'client' => 'required|uuid',
            'secret' => 'required',
            'super' => 'required'
        ]);
        if ($validate->fails()) {
            return [
                'error' => 1,
                'msg' => $validate->errors()
            ];
        }
        return $this->curd
            ->addModel('real_time_client')
            ->afterHook(function () {
                $this->authRefresh();
                return true;
            })
            ->result();
    }

    public function edit(): array
    {
        $validate = $this->curd->editValidation([
            'group' => 'required',
            'secret' => 'required',
            'super' => 'required'
        ]);
        if ($validate->fails()) {
            return [
                'error' => 1,
                'msg' => $validate->errors()
            ];
        }
        return $this->curd
            ->editModel('real_time_client')
            ->afterHook(function () {
                $this->authRefresh();
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
            ->deleteModel('real_time_client')
            ->afterHook(function () {
                $this->authRefresh();
                return true;
            })
            ->result();
    }

    /**
     * 获取绑定
     * @return array
     */
    public function getBinding(): array
    {
        $body = $this->request->post();
        $validator = $this->validation->make($body, [
            'client' => 'required|uuid',
        ]);
        if ($validator->fails()) {
            return [
                'error' => 1,
                'msg' => $validator->errors()
            ];
        }

        $query = Db::table('real_time_binding')
            ->where('client', '=', $body['client'])
            ->get();

        return [
            'error' => 0,
            'data' => $query->toArray()
        ];
    }

    /**
     * 绑定更新
     * @return array
     */
    public function putBinding(): array
    {
        $body = $this->request->post();
        $validator = $this->validation->make($body, [
            'client' => 'required|uuid',
            'topic' => 'required',
            'policy' => 'required'
        ]);
        if ($validator->fails()) {
            return [
                'error' => 1,
                'msg' => $validator->errors()
            ];
        }
        $query = Db::table('real_time_client')
            ->where('client', '=', $body['client'])
            ->first();

        if (empty($query)) {
            return [
                'error' => 1,
                'msg' => 'client not exists'
            ];
        }
        if ($query->super === 1) {
            return [
                'error' => 1,
                'msg' => 'client is super'
            ];
        }
        $result = Db::table('real_time_binding')
            ->insert([
                'client' => $body['client'],
                'topic' => $body['topic'],
                'policy' => $body['policy']
            ]);
        $this->authRefresh();
        return $result ? [
            'error' => 0,
            'msg' => 'ok'
        ] : [
            'error' => 1,
            'msg' => 'failed'
        ];
    }

    /**
     * 绑定移除
     * @return array
     */
    public function deleteBinding(): array
    {
        $body = $this->request->post();
        $validator = $this->validation->make($body, [
            'id' => 'required',
        ]);
        if ($validator->fails()) {
            return [
                'error' => 1,
                'msg' => $validator->errors()
            ];
        }
        $result = Db::table('real_time_binding')
            ->where('id', '=', $body['id'])
            ->delete();

        $this->authRefresh();
        return $result !== 0 ? [
            'error' => 0,
            'msg' => 'ok'
        ] : [
            'error' => 1,
            'msg' => 'failed'
        ];
    }

    /**
     * 更新授权缓存
     */
    private function authRefresh(): void
    {
        $this->eventDispatcher->dispatch(new RealTimeRefresh());
    }
}
