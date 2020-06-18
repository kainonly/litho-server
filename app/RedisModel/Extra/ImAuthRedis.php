<?php
declare (strict_types=1);

namespace App\RedisModel\Extra;

use Hyperf\DbConnection\Db;

class ImAuthRedis extends RedisModel
{
    protected string $key = 'im-auth';

    /**
     * 刷新授权缓存
     */
    public function refresh(): void
    {
        $query = Db::table('im_key')
            ->get();

        if ($query->isEmpty()) {
            return;
        }

        $lists = [];
        foreach ($query->toArray() as $value) {
            $lists[$value->key] = $value->secret;
        }
        $multi = $this->redis->multi();
        $multi->del($this->key);
        $multi->hMSet($this->key, $lists);
        $multi->exec();
    }
}