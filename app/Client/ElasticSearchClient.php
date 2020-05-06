<?php
declare(strict_types=1);

namespace App\Client;

use Elasticsearch\Client;
use Hyperf\Contract\ConfigInterface;
use Hyperf\Elasticsearch\ClientBuilderFactory;
use stdClass;

class ElasticSearchClient
{
    /**
     * @var stdClass
     */
    private stdClass $client;

    public function __construct(
        ClientBuilderFactory $factory,
        ConfigInterface $config
    )
    {
        $builder = $factory->create();
        $this->client = new stdClass();
        foreach ($config->get('elasticsearch') as $key => $value) {
            $this->client->$key = $builder->setHosts($value['hosts'])->build();
        }
    }

    /**
     * 获取客户端
     * @param string $key
     * @return Client
     */
    public function client(string $key = 'default'): Client
    {
        return $this->client->$key;
    }
}