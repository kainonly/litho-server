<?php
declare (strict_types=1);

namespace App\HttpClient;

use Hyperf\Guzzle\ClientFactory;
use RabbitMQ\API\RabbitMQ;
use stdClass;

class RabbitMQService
{
    /**
     * @var RabbitMQ
     */
    private RabbitMQ $master;
    /**
     * @var string
     */
    private string $masterNode;
    /**
     * @var string
     */
    private string $masterAMQP;
    /**
     * @var stdClass
     */
    private stdClass $broker;
    /**
     * @var stdClass
     */
    private stdClass $brokerNode;
    /**
     * @var stdClass
     */
    private stdClass $brokerAMQP;

    public function __construct(ClientFactory $clientFactory)
    {
        $option = config('rabbitmq');
        $this->masterNode = $option['master']['node'];
        $this->masterAMQP = $option['master']['amqp'];
        $this->master = new RabbitMQ($clientFactory->create([
            'base_uri' => $option['master']['uri'],
            'auth' => [$option['master']['user'], $option['master']['pass']],
            'timeout' => 5.0,
            'version' => (float)$option['master']['version']
        ]));
        $this->brokerNode = new stdClass();
        $this->brokerAMQP = new stdClass();
        $this->broker = new stdClass();
        foreach ($option['broker'] as $key => $value) {
            $this->brokerNode->$key = $value['node'];
            $this->brokerAMQP->$key = $value['amqp'];
            $this->broker->$key = new RabbitMQ($clientFactory->create([
                'base_uri' => $value['uri'],
                'auth' => [$value['user'], $value['pass']],
                'timeout' => 5.0,
                'version' => (float)$value['version']
            ]));
        }
    }

    /**
     * @return RabbitMQ
     */
    public function getMaster(): RabbitMQ
    {
        return $this->master;
    }

    /**
     * @return string
     */
    public function getMasterNode(): string
    {
        return $this->masterNode;
    }

    /**
     * @return string
     */
    public function getMasterAMQP(): string
    {
        return $this->masterAMQP;
    }

    /**
     * @param string $label
     * @return RabbitMQ
     */
    public function getBroker(string $label): RabbitMQ
    {
        return $this->broker->$label;
    }

    /**
     * @param string $label
     * @return string
     */
    public function getBrokerNode(string $label): string
    {
        return $this->brokerNode->$label;
    }

    /**
     * @param string $label
     * @return string
     */
    public function getBrokerAMQP(string $label): string
    {
        return $this->brokerAMQP->$label;
    }
}