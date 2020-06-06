<?php
declare (strict_types=1);

namespace App\RedisModel\Common;

use Hyperf\DbConnection\Db;
use Hyperf\Extra\Common\RedisModel;

class RealTimeAuthRedis extends RedisModel
{
    protected string $key = 'common:real-time-auth';

    public function refresh(): void
    {
        $query = Db::table('real_time_client')
            ->get();

        if ($query->isEmpty()) {
            return;
        }

        $lists = [];
        foreach ($query->toArray() as $value) {
            $lists[$value->client] = $value->secret;
        }
        $this->redis->hMSet($this->key, $lists);
    }
}