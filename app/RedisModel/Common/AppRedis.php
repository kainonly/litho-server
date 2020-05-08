<?php
declare (strict_types=1);

namespace App\RedisModel\Common;

use Hyperf\DbConnection\Db;
use Hyperf\Extra\Common\RedisModel;

class AppRedis extends RedisModel
{
    protected string $key = 'common:app';

    public function clear(): void
    {
        $this->redis->del($this->key);
    }

    public function get(string $appid): array
    {
        if (!$this->redis->exists($this->key)) {
            $this->update();
        }
        $raws = $this->redis->hGet($this->key, $appid);
        return !empty($raws) ? json_decode($raws, true) : [];
    }

    public function update(): void
    {
        $query = Db::table('app')
            ->where('status', '=', 1)
            ->get(['name', 'appid', 'secret', 'entry', 'expires']);

        if ($query->isEmpty()) {
            return;
        }

        $lists = [];
        foreach ($query->toArray() as $value) {
            $lists[$value->appid] = json_encode($value);
        }
        $this->redis->hMSet($this->key, $lists);
    }
}