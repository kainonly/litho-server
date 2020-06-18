<?php
declare(strict_types=1);

namespace App\RedisModel\Extra;

use App\Client\ExtraRedis;
use Hyperf\Redis\Redis;
use Psr\Container\ContainerInterface;

abstract class RedisModel
{
    protected string $key;
    protected ContainerInterface $container;
    /**
     * @var Redis|\Redis
     */
    protected Redis $redis;

    /**
     * RedisModel constructor.
     * @param ContainerInterface $container
     */
    public function __construct(ContainerInterface $container)
    {
        $this->container = $container;
        $this->redis = $container->get(ExtraRedis::class);
    }
}
