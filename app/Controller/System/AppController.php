<?php
declare(strict_types=1);

namespace App\Controller\System;

use App\RedisModel\Common\AppRedis;
use Hyperf\DbConnection\Db;
use Hyperf\Di\Annotation\Inject;
use Hyperf\Utils\Str;

class AppController extends BaseController
{
    /**
     * @Inject()
     * @var AppRedis
     */
    private AppRedis $appRedis;

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
            ->originListsModel('app')
            ->setField(['id', 'name', 'appid'])
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
            ->listsModel('app')
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
            ->getModel('app')
            ->result();
    }

    public function add(): array
    {
        $body = $this->request->post();
        $validate = $this->curd->addValidation();
        if ($validate->fails()) {
            return [
                'error' => 1,
                'msg' => $validate->errors()
            ];
        }
        if (!empty($body['expires']) && $body['expires'] !== 0) {
            $body['expires'] = strtotime($body['expires']);
        }
        return $this->curd
            ->addModel('app')
            ->afterHook(function () {
                $this->clearRedis();
                return true;
            })
            ->result();
    }

    public function edit(): array
    {
        $body = $this->request->post();
        $validate = $this->curd->editValidation();
        if ($validate->fails()) {
            return [
                'error' => 1,
                'msg' => $validate->errors()
            ];
        }
        if (!empty($body['expires']) && $body['expires'] !== 0) {
            $body['expires'] = strtotime($body['expires']);
        }
        return $this->curd
            ->editModel('app')
            ->afterHook(function () {
                $this->clearRedis();
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
            ->deleteModel('app')
            ->afterHook(function () {
                $this->clearRedis();
                return true;
            })
            ->result();
    }

    /**
     * 清除缓存
     */
    private function clearRedis(): void
    {
        $this->appRedis->clear();
    }

    /**
     * 随机生成appid
     * @return array
     */
    public function randomAppid(): array
    {
        return [
            'error' => 0,
            'data' => Str::random()
        ];
    }

    /**
     * 随机生成secret
     * @return array
     */
    public function randomSecret(): array
    {
        return [
            'error' => 0,
            'data' => sha1(Str::random())
        ];
    }

    /**
     * 验证应用名称
     * @return array
     */
    public function validedName(): array
    {
        $body = $this->request->post();
        if (empty($body['name'])) {
            return [
                'error' => 1,
                'msg' => 'error:has_name'
            ];
        }

        $result = Db::table('app')
            ->where('name', '=', $body['name'])
            ->count();

        return [
            'error' => 0,
            'data' => !empty($result)
        ];
    }

    /**
     * 验证APPID是否存在
     * @return array
     */
    public function validedAppid(): array
    {
        $body = $this->request->post();
        if (empty($body['appid'])) {
            return [
                'error' => 1,
                'msg' => 'error:has_name'
            ];
        }

        $result = Db::table('app')
            ->where('appid', '=', $body['appid'])
            ->count();

        return [
            'error' => 0,
            'data' => !empty($result)
        ];
    }
}