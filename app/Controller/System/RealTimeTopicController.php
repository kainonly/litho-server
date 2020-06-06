<?php
declare(strict_types=1);

namespace App\Controller\System;

use App\Event\RealTimeRefresh;
use Hyperf\DbConnection\Db;
use Hyperf\Di\Annotation\Inject;
use Psr\EventDispatcher\EventDispatcherInterface;

class RealTimeTopicController extends BaseController
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
            ->originListsModel('real_time_topic')
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
            ->listsModel('real_time_topic')
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
            ->getModel('real_time_topic')
            ->result();
    }

    public function add(): array
    {
        $validate = $this->curd->addValidation([
            'topic' => 'required'
        ]);
        if ($validate->fails()) {
            return [
                'error' => 1,
                'msg' => $validate->errors()
            ];
        }
        return $this->curd
            ->addModel('real_time_topic')
            ->afterHook(function () {
                $this->authRefresh();
                return true;
            })
            ->result();
    }

    public function edit(): array
    {
        $validate = $this->curd->editValidation([
            'topic' => 'required_if:switch,false'
        ]);
        if ($validate->fails()) {
            return [
                'error' => 1,
                'msg' => $validate->errors()
            ];
        }
        return $this->curd
            ->editModel('real_time_topic')
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
            ->deleteModel('real_time_topic')
            ->afterHook(function () {
                $this->authRefresh();
                return true;
            })
            ->result();
    }

    /**
     * Exists Topic Key
     * @return array
     */
    public function validedTopic(): array
    {
        $body = $this->request->post();
        if (empty($body['topic'])) {
            return [
                'error' => 1,
                'msg' => 'require key'
            ];
        }

        $exists = Db::table('real_time_topic')
            ->where('topic', '=', $body['topic'])
            ->exists();

        return [
            'error' => 0,
            'data' => $exists
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
