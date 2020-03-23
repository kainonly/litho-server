<?php
declare (strict_types=1);

namespace HyperfTest\Microservice;

use AMQPSubscriber\AllResponse;
use AMQPSubscriber\DeleteParameter;
use AMQPSubscriber\GetParameter;
use AMQPSubscriber\GetResponse;
use AMQPSubscriber\ListsParameter;
use AMQPSubscriber\ListsResponse;
use AMQPSubscriber\NoParameter;
use AMQPSubscriber\PutParameter;
use AMQPSubscriber\Response;
use AMQPSubscriber\RouterClient;
use Exception;
use Grpc\ChannelCredentials;
use Hyperf\Utils\Filesystem\Filesystem;
use PHPStan\Testing\TestCase;

class AMQPServiceTest extends TestCase
{
    private RouterClient $client;
    private array $config;

    public function __construct($name = null, array $data = [], $dataName = '')
    {
        parent::__construct($name, $data, $dataName);
        try {
            $this->client = new RouterClient('127.0.0.1:6000', [
                'credentials' => ChannelCredentials::createInsecure(),
            ]);
            $filesystem = new Filesystem();
            $context = $filesystem->get('../config/amqp/config.yml');
            $this->config = yaml_parse($context);
        } catch (Exception $e) {
            $this->expectException($e->getMessage());
        }
    }

    public function testPut()
    {
        $param = new PutParameter();
        $param->setIdentity('a1');
        $param->setQueue('test');
        $param->setUrl('http://10.0.75.1:3000');
        $param->setSecret('123');
        /**
         * @var Response $response
         */
        list($response, $status) = $this->client->Put($param)->wait();
        $this->assertEquals($response->getError(), 0, $response->getMsg());
    }

    public function testAll()
    {
        $param = new NoParameter();
        /**
         * @var AllResponse $response
         */
        list($response, $status) = $this->client->All($param)->wait();
        $this->assertEquals($response->getError(), 0, $response->getMsg());
        foreach ($response->getData()->getIterator() as $value) {
            var_dump($value);
        }
    }

    public function testGet()
    {
        $param = new GetParameter();
        $param->setIdentity('a1');
        /**
         * @var GetResponse $response
         */
        list($response, $status) = $this->client->Get($param)->wait();
        $this->assertEquals($response->getError(), 0, $response->getMsg());
        var_dump($response->getData()->serializeToJsonString());
    }

    public function testLists()
    {
        $param = new ListsParameter();
        $param->setIdentity(['a1', 'a2']);
        /**
         * @var ListsResponse $response
         */
        list($response, $status) = $this->client->Lists($param)->wait();
        $this->assertEquals($response->getError(), 0, $response->getMsg());
        foreach ($response->getData()->getIterator() as $value) {
            var_dump($value->serializeToJsonString());
        }
    }

    public function testDelete()
    {
        $param = new DeleteParameter();
        $param->setIdentity('a1');
        /**
         * @var Response $response
         */
        list($response, $status) = $this->client->Delete($param)->wait();
        $this->assertEquals($response->getError(), 0, $response->getMsg());
    }
}