<?php
declare (strict_types=1);

namespace HyperfTest\Microservice;

use Exception;
use Grpc\ChannelCredentials;
use Hyperf\Utils\Filesystem\Filesystem;
use PHPStan\Testing\TestCase;
use SSHMicroservice\AllResponse;
use SSHMicroservice\DeleteParameter;
use SSHMicroservice\ExecParameter;
use SSHMicroservice\ExecResponse;
use SSHMicroservice\GetParameter;
use SSHMicroservice\GetResponse;
use SSHMicroservice\ListsParameter;
use SSHMicroservice\ListsResponse;
use SSHMicroservice\NoParameter;
use SSHMicroservice\PutParameter;
use SSHMicroservice\Response;
use SSHMicroservice\RouterClient;
use SSHMicroservice\TunnelOption;
use SSHMicroservice\TunnelsParameter;

class SSHServiceTest extends TestCase
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
            $context = $filesystem->get('../config/debug.yml');
            $this->config = yaml_parse($context);
            $key = $filesystem->get('../config/key-1.pem');
            $this->config[0]['private_key'] = base64_encode($key);
        } catch (Exception $e) {
            $this->expectException($e->getMessage());
        }
    }

    public function testPut()
    {
        $param = new PutParameter();
        $param->setIdentity('test');
        $param->setHost($this->config[0]['host']);
        $param->setPort($this->config[0]['port']);
        $param->setUsername($this->config[0]['username']);
        $param->setPassword($this->config[0]['password'] ?? '');
        $param->setPrivateKey($this->config[0]['private_key']);
        $param->setPassphrase($this->config[0]['passphrase'] ?? '');
        /**
         * @var Response $response
         */
        list($response, $status) = $this->client->Put($param)->wait();
        $this->assertEquals($response->getError(), 0, $response->getMsg());
    }

    public function testExec()
    {
        $param = new ExecParameter();
        $param->setIdentity('test');
        $param->setBash('uptime');
        /**
         * @var ExecResponse $response
         */
        list($response, $status) = $this->client->Exec($param)->wait();
        $this->assertEquals($response->getError(), 0, $response->getMsg());
        var_dump($response->getData());
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

    private function setTunnelOption(
        string $src_ip,
        int $src_port,
        string $dst_ip,
        int $dst_port
    ): TunnelOption
    {
        $option = new TunnelOption();
        $option->setSrcIp($src_ip);
        $option->setSrcPort($src_port);
        $option->setDstIp($dst_ip);
        $option->setDstPort($dst_port);
        return $option;
    }

    public function testTunnels()
    {
        $param = new TunnelsParameter();
        $param->setIdentity('test');
        $param->setTunnels([
            $this->setTunnelOption('127.0.0.1', 3306, '127.0.0.1', 3306),
            $this->setTunnelOption('127.0.0.1', 9200, '127.0.0.1', 9200)
        ]);
        /**
         * @var Response $response
         */
        list($response, $status) = $this->client->Tunnels($param)->wait();
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