<?php
declare (strict_types=1);

namespace App\Listener;

use App\Event\ImRefresh;
use App\RedisModel\Common\ImKeyRedis;
use App\RedisModel\Extra\ImAclRedis;
use App\RedisModel\Extra\ImAuthRedis;
use App\RedisModel\Extra\ImSuperRedis;
use App\RedisModel\Common\ImSuperRedis as CommonImSuperRedis;
use App\RedisModel\Common\ImAclRedis as CommonImAclRedis;
use Hyperf\Di\Annotation\Inject;
use Hyperf\Event\Contract\ListenerInterface;

class ImListener implements ListenerInterface
{
    /**
     * @Inject()
     * @var ImKeyRedis
     */
    private ImKeyRedis $imKeyRedis;
    /**
     * @Inject()
     * @var ImAuthRedis
     */
    private ImAuthRedis $imAuthRedis;
    /**
     * @Inject()
     * @var CommonImSuperRedis
     */
    private CommonImSuperRedis $cImSuperRedis;
    /**
     * @Inject()
     * @var ImSuperRedis
     */
    private ImSuperRedis $imSuperRedis;
    /**
     * @Inject()
     * @var CommonImAclRedis
     */
    private CommonImAclRedis $cImAclRedis;
    /**
     * @Inject()
     * @var ImAclRedis
     */
    private ImAclRedis $imAclRedis;

    public function listen(): array
    {
        return [
            ImRefresh::class,
        ];
    }

    public function process(object $event): void
    {
        if ($event instanceof ImRefresh) {
            $this->imKeyRedis->clear();
            $this->imAuthRedis->refresh();
            $this->imSuperRedis->refresh();
            $this->cImSuperRedis->refresh();
            $this->imAclRedis->refresh();
            $this->cImAclRedis->refresh();
        }
    }
}