<?php
declare (strict_types=1);

namespace App\GrpcClient;

use Hyperf\Contract\ConfigInterface;
use Hyperf\Pool\SimplePool\PoolFactory;
use Psr\Container\ContainerInterface;

class AMQPServiceFactory
{
    public function __invoke(ContainerInterface $container)
    {
        /**
         * @var ConfigInterface $config
         */
        $config = $container->get(ConfigInterface::class);
        $option = $config->get('microservice.amqp');
        /**
         * @var PoolFactory $factory
         */
        $factory = $container->get(PoolFactory::class);
        $pool = $factory->get('amqp-service-pool', function () use ($option) {
            return new AMQPService($option['host'], [
                'credentials' => null,
            ]);
        }, [
            'max_connections' => $option['max_connections']
        ]);
        return $pool->get()->getConnection();
    }
}