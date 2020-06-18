<?php
declare (strict_types=1);

namespace App\RedisModel\Common;

use Hyperf\DbConnection\Db;
use Hyperf\Extra\Common\RedisModel;

class ImKeyRedis extends RedisModel
{
    protected string $key = 'common:im-key';

    public function clear(): void
    {
        $this->redis->del($this->key);
    }

    public function get(string $key): array
    {
        if (!$this->redis->exists($this->key)) {
            $this->update();
        }
        $raws = $this->redis->hGet($this->key, $key);
        return !empty($raws) ? json_decode($raws, true) : [];
    }

    public function update(): void
    {
        $query = Db::table('im_key')
            ->get(['group', 'key', 'secret', 'super']);

        if ($query->isEmpty()) {
            return;
        }

        $lists = [];
        foreach ($query->toArray() as $value) {
            $lists[$value->key] = json_encode($value);
        }
        $this->redis->hMSet($this->key, $lists);
    }
}