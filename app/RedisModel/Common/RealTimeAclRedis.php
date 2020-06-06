<?php
declare (strict_types=1);

namespace App\RedisModel\Common;

use Hyperf\DbConnection\Db;
use Hyperf\Extra\Common\RedisModel;
use Psr\Container\ContainerInterface;

class RealTimeAclRedis extends RedisModel
{
    protected string $key = 'common:real-time-acl:';

    public function refresh(): void
    {
        $exists = $this->redis->keys($this->key . '*');
        if (!empty($exists)) {
            $this->redis->del($exists);
        }
        $query = Db::table('real_time_binding')
            ->get();

        if ($query->isEmpty()) {
            return;
        }

        $clients = [];
        foreach ($query->toArray() as $value) {
            if ($value->policy === 0) {
                $clients[$value->client][] = $value->topic;
            }
            if ($value->policy === 1) {
                $clients[$value->client][] = $value->topic . '/' . $value->client;
            }
        }
        foreach ($clients as $client => $acl) {
            $this->redis->sAdd($this->key . $client, ...$acl);
        }
    }
}