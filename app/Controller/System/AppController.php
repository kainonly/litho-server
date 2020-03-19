<?php
declare(strict_types=1);

namespace App\Controller\System;

use Hyperf\Curd\Common\AddModel;
use Hyperf\Curd\Common\DeleteModel;
use Hyperf\Curd\Common\EditModel;
use Hyperf\Curd\Common\GetModel;
use Hyperf\Curd\Common\ListsModel;
use Hyperf\Curd\Common\OriginListsModel;
use Hyperf\Curd\Lifecycle\AddAfterHooks;
use Hyperf\Curd\Lifecycle\AddBeforeHooks;
use Hyperf\Curd\Lifecycle\DeleteAfterHooks;
use Hyperf\Curd\Lifecycle\EditAfterHooks;
use Hyperf\Curd\Lifecycle\EditBeforeHooks;

class AppController extends BaseController
    implements AddBeforeHooks, AddAfterHooks, EditBeforeHooks, EditAfterHooks, DeleteAfterHooks
{
    use OriginListsModel, ListsModel, AddModel, GetModel, EditModel, DeleteModel;
    protected string $model = 'app';
    protected array $origin_lists_field = ['id', 'name', 'appid'];
    private string $appid;

    /**
     * @inheritDoc
     */
    public function addBeforeHooks(): bool
    {
        if (!empty($this->post['expires']) && $this->post['expires'] !== 0) {
            $this->post['expires'] = strtotime($this->post['expires']);
        }
        return true;
    }

    /**
     * @inheritDoc
     */
    public function addAfterHooks(int $id): bool
    {
    }

    /**
     * @inheritDoc
     */
    public function editBeforeHooks(): bool
    {
        // TODO: Implement editBeforeHooks() method.
    }

    /**
     * @inheritDoc
     */
    public function editAfterHooks(): bool
    {
        // TODO: Implement editAfterHooks() method.
    }

    /**
     * @inheritDoc
     */
    public function deleteAfterHooks(): bool
    {
        // TODO: Implement deleteAfterHooks() method.
    }
}