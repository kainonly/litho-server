<?php
declare (strict_types=1);

namespace HyperfTest\Microservice;

use Exception;
use Grpc\ChannelCredentials;
use Hyperf\Utils\Filesystem\Filesystem;
use PHPStan\Testing\TestCase;
use ScheduleMicroservice\AllResponse;
use ScheduleMicroservice\DeleteParameter;
use ScheduleMicroservice\EntryOption;
use ScheduleMicroservice\GetParameter;
use ScheduleMicroservice\GetResponse;
use ScheduleMicroservice\ListsParameter;
use ScheduleMicroservice\ListsResponse;
use ScheduleMicroservice\NoParameter;
use ScheduleMicroservice\PutParameter;
use ScheduleMicroservice\Response;
use ScheduleMicroservice\RouterClient;
use ScheduleMicroservice\RunningParameter;

class ScheduleServiceTest extends TestCase
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
            $context = $filesystem->get('../config/schedule/config.yml');
            $this->config = yaml_parse($context);
        } catch (Exception $e) {
            $this->expectException($e->getMessage());
        }
    }

    private function entryOption(string $cron_time, string $url, array $headers, array $body): EntryOption
    {
        $option = new EntryOption();
        $option->setCronTime($cron_time);
        $option->setUrl($url);
        $option->setHeaders(json_encode($headers));
        $option->setBody(json_encode($body));
        return $option;
    }

    public function testPut()
    {
        $param = new PutParameter();
        $param->setIdentity('test');
        $param->setTimeZone('Asia/Shanghai');
        $param->setStart(true);
        $param->setEntries([
            'task1' => $this->entryOption(
                '*/5 * * * * *',
                'http://10.0.75.1:3000',
                ['x-token' => 'abc'],
                ['name' => 'task1']
            ),
            'task2' => $this->entryOption(
                '*/10 * * * * *',
                'http://10.0.75.1:3000',
                ['x-token' => 'abc'],
                ['name' => 'task2']
            )
        ]);
        /**
         * @var Response $response
         */
        list($response, $status) = $this->client->Put($param)->wait();
        $this->assertEquals($response->getError(), 0, $response->getMsg());
    }

    public function testGet()
    {
        $param = new GetParameter();
        $param->setIdentity('test');
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
        $param->setIdentity(['test']);
        /**
         * @var ListsResponse $response
         */
        list($response, $status) = $this->client->Lists($param)->wait();
        $this->assertEquals($response->getError(), 0, $response->getMsg());
        foreach ($response->getData()->getIterator() as $value) {
            var_dump($value->serializeToJsonString());
        }
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

    public function testRuning()
    {
        $param = new RunningParameter();
        $param->setIdentity('test');
        $param->setRunning(false);
        /**
         * @var Response $response
         */
        list($response, $status) = $this->client->Running($param)->wait();
        $this->assertEquals($response->getError(), 0, $response->getMsg());
    }

    public function testDelete()
    {
        $param = new DeleteParameter();
        $param->setIdentity('test');
        /**
         * @var Response $response
         */
        list($response, $status) = $this->client->Delete($param)->wait();
        $this->assertEquals($response->getError(), 0, $response->getMsg());
    }
}