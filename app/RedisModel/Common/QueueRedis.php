<?php
declare (strict_types=1);

namespace App\RedisModel\Common;

use Hyperf\DbConnection\Db;
use Hyperf\Extra\Common\RedisModel;

class QueueRedis extends RedisModel
{
    protected string $key = 'common:queue';

    public function clear(): void
    {
        $this->redis->del($this->key);
    }

    public function get(string $queue): array
    {
        if (!$this->redis->exists($this->key)) {
            $this->update();
        }
        $raws = $this->redis->hGet($this->key, $queue);
        return !empty($raws) ? json_decode($raws, true) : [];
    }

    public function update(): void
    {
        $query = Db::table('queue')->get();

        if ($query->isEmpty()) {
            return;
        }

        $lists = [];
        foreach ($query->toArray() as $value) {
            $lists[$value->queue] = json_encode([
                'group' => $value->group,
                'sizeof' => $value->sizeof * $value->sizeof_unit,
                'ttl' => $value->ttl * $value->ttl_unit,
                'limit' => $value->limit * $value->limit_unit,
                'lazy' => $value->lazy,
                'rewind' => $value->rewind === 1 ? $value->rewind_ttl * $value->rewind_ttl_unit : 0
            ]);
        }
        $this->redis->hMSet($this->key, $lists);
    }
}