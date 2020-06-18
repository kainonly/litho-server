<?php
declare(strict_types=1);

namespace App\Client;

use Hyperf\Redis\Redis;

class ExtraRedis extends Redis
{
    protected $poolName = 'extra';
}