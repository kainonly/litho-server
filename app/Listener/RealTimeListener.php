<?php
declare (strict_types=1);

namespace App\Listener;

use App\Event\RealTimeRefresh;
use App\RedisModel\Common\RealTimeAclRedis;
use App\RedisModel\Common\RealTimeAuthRedis;
use App\RedisModel\Common\RealTimeSuperRedis;
use Hyperf\Di\Annotation\Inject;
use Hyperf\Event\Contract\ListenerInterface;

class RealTimeListener implements ListenerInterface
{
    /**
     * @Inject()
     * @var RealTimeAuthRedis
     */
    private RealTimeAuthRedis $realTimeAuthRedis;
    /**
     * @Inject()
     * @var RealTimeSuperRedis
     */
    private RealTimeSuperRedis $realTimeSuperRedis;
    /**
     * @Inject()
     * @var RealTimeAclRedis
     */
    private RealTimeAclRedis $realTimeAclRedis;

    public function listen(): array
    {
        return [
            RealTimeRefresh::class,
        ];
    }

    public function process(object $event)
    {
        if ($event instanceof RealTimeRefresh) {
            $this->realTimeAuthRedis->refresh();
            $this->realTimeSuperRedis->refresh();
            $this->realTimeAclRedis->refresh();
        }
    }
}