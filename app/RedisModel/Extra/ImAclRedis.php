<?php
declare (strict_types=1);

namespace App\RedisModel\Extra;

use Hyperf\DbConnection\Db;

class ImAclRedis extends RedisModel
{
    protected string $key = 'im-acl:';

    /**
     * 刷新权限缓存
     */
    public function refresh(): void
    {
        $exists = $this->redis->keys($this->key . '*');
        if (!empty($exists)) {
            $this->redis->del($exists);
        }
        $query = Db::table('im_binding')
            ->get();

        if ($query->isEmpty()) {
            return;
        }

        $keys = [];
        foreach ($query->toArray() as $value) {
            if ($value->policy === 0) {
                $keys[$value->key][] = $value->topic;
            }
            if ($value->policy === 1) {
                $keys[$value->key][] = $value->topic . '/' . $value->key;
            }
        }
        $multi = $this->redis->multi();
        foreach ($keys as $key => $acl) {
            $multi->del($this->key);
            $multi->sAdd($this->key . $key, ...$acl);
        }
        $multi->exec();
    }
}