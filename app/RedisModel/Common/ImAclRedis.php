<?php
declare (strict_types=1);

namespace App\RedisModel\Common;

use Hyperf\DbConnection\Db;
use Hyperf\Extra\Common\RedisModel;

class ImAclRedis extends RedisModel
{
    protected string $key = 'common:im-acl:';

    /**
     * 确认是否存在
     * @param string $key 通行码
     * @param string $topic 主题名称
     * @return bool
     */
    public function is(string $key, string $topic): bool
    {
        return $this->redis->sIsMember($this->key . $key, $topic);
    }

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