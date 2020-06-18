<?php
declare (strict_types=1);

namespace App\RedisModel\Common;

use Hyperf\DbConnection\Db;
use Hyperf\Extra\Common\RedisModel;

class ImSuperRedis extends RedisModel
{
    protected string $key = 'common:im-super';

    /**
     * 确认是否存在
     * @param string $key 通行码
     * @return bool
     */
    public function is(string $key): bool
    {
        return $this->redis->sIsMember($this->key, $key);
    }

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