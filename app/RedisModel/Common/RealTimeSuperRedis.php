<?php
declare (strict_types=1);

namespace App\RedisModel\Common;

use Hyperf\DbConnection\Db;
use Hyperf\Extra\Common\RedisModel;

class RealTimeSuperRedis extends RedisModel
{
    protected string $key = 'common:real-time-super';

    public function refresh(): void
    {
        $query = Db::table('real_time_client')
            ->where('super', '=', 1)
            ->get();

        if ($query->isEmpty()) {
            return;
        }

        $lists = [];
        foreach ($query->toArray() as $value) {
            $lists[] = $value->client;
        }
        $this->redis->sAdd($this->key, ...$lists);
    }
}