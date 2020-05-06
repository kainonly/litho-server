<?php
declare(strict_types=1);

namespace App\Client;

use Hyperf\AMQPClient\AMQPClientInterface;
use Hyperf\Di\Annotation\Inject;
use PhpAmqpLib\Connection\AMQPStreamConnection;
use stdClass;

class AMQPClient
{
    /**
     * @Inject()
     * @var AMQPClientInterface
     */
    private AMQPClientInterface $amqp;

    /**
     * @var stdClass
     */
    private stdClass $connection;

    /**
     * AMQPClient constructor.
     */
    public function __construct()
    {
        $this->connection = new stdClass();
    }

    /**
     * @param string $key
     * @return AMQPStreamConnection
     */
    public function client(string $key): AMQPStreamConnection
    {
        if (empty($this->connection->$key)) {
            $this->connection->$key = $this->amqp
                ->client('default')
                ->getAMQPStreamConnection();
        }
        return $this->connection->$key;
    }

    public function delete(string $key): void
    {
        unset($this->connection->$key);
    }
}