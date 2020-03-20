<?php
declare (strict_types=1);

namespace HyperfTest\Microservice;

use Grpc\ChannelCredentials;
use PHPStan\Testing\TestCase;
use SSHMicroservice\AllResponse;
use SSHMicroservice\NoParameter;
use SSHMicroservice\RouterClient;

class SSHServiceTest extends TestCase
{
    private RouterClient $client;

    public function __construct($name = null, array $data = [], $dataName = '')
    {
        parent::__construct($name, $data, $dataName);
        $this->client = new RouterClient('127.0.0.1:6000', [
            'credentials' => ChannelCredentials::createInsecure(),
        ]);
    }

    public function testAll()
    {
        $param = new NoParameter();
        /**
         * @var AllResponse $response
         */
        list($response, $status) = $this->client->All($param)->wait();
        $this->assertEquals($response->getError(), 0);
        foreach ($response->getData()->getIterator() as $value) {
            var_dump($value);
        }
    }
}