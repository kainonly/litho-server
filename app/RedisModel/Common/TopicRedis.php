<?php
declare (strict_types=1);

namespace App\RedisModel\Common;

use Hyperf\DbConnection\Db;
use Hyperf\Extra\Common\RedisModel;

class TopicRedis extends RedisModel
{
    protected string $key = 'common:topic';

    public function clear(): void
    {
        $this->redis->del($this->key);
    }

    public function get(string $topic): array
    {
        if (!$this->redis->exists($this->key)) {
            $this->update();
        }
        $raws = $this->redis->hGet($this->key, $topic);
        return !empty($raws) ? json_decode($raws, true) : [];
    }

    public function update(): void
    {
        $query = Db::table('topic')->get();

        if ($query->isEmpty()) {
            return;
        }

        $lists = [];
        foreach ($query->toArray() as $value) {
            $lists[$value->topic] = json_encode([
                'group' => $value->group,
                'sizeof' => $value->sizeof * $value->sizeof_unit
            ]);
        }
        $this->redis->hMSet($this->key, $lists);
    }
}