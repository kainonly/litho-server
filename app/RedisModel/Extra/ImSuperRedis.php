<?php
declare (strict_types=1);

namespace App\RedisModel\Extra;

use Hyperf\DbConnection\Db;

class ImSuperRedis extends RedisModel
{
    protected string $key = 'im-super';

    /**
     * 刷新超级用户
     */
    public function refresh(): void
    {
        $query = Db::table('im_key')
            ->where('super', '=', 1)
            ->get();

        if ($query->isEmpty()) {
            return;
        }

        $lists = [];
        foreach ($query->toArray() as $value) {
            $lists[] = $value->key;
        }
        $multi = $this->redis->multi();
        $multi->del($this->key);
        $multi->sAdd($this->key, ...$lists);
        $multi->exec();
    }
}