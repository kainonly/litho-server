<?php
declare(strict_types=1);

namespace App\Controller\System;

use App\Client\ElasticSearchClient;
use App\GrpcClient\AMQPSubscriberServiceInterface;
use App\HttpClient\RabbitMQService;
use App\RedisModel\Common\QueueRedis;
use Exception;
use Hyperf\DbConnection\Db;
use Hyperf\Di\Annotation\Inject;
use Hyperf\Utils\Context;
use RabbitMQ\API\Common\ExchangeOption;
use RabbitMQ\API\Common\MessageOption;
use RabbitMQ\API\Common\PublishOption;
use RabbitMQ\API\Common\QueueOption;

class QueueController extends BaseController
{
    /**
     * @Inject()
     * @var RabbitMQService
     */
    private RabbitMQService $rabbit;
    /**
     * @Inject()
     * @var QueueRedis
     */
    private QueueRedis $queueRedis;
    /**
     * @Inject()
     * @var ElasticSearchClient
     */
    private ElasticSearchClient $es;
    /**
     * @Inject()
     * @var AMQPSubscriberServiceInterface
     */
    private AMQPSubscriberServiceInterface $amqpSubscriberService;

    public function originLists(): array
    {
        $validate = $this->curd->originListsValidation();
        if ($validate->fails()) {
            return [
                'error' => 1,
                'msg' => $validate->errors()
            ];
        }

        return $this->curd
            ->originListsModel('queue')
            ->setOrder('create_time', 'desc')
            ->result();
    }

    public function lists(): array
    {
        $validate = $this->curd->listsValidation();
        if ($validate->fails()) {
            return [
                'error' => 1,
                'msg' => $validate->errors()
            ];
        }

        return $this->curd
            ->listsModel('queue')
            ->setOrder('create_time', 'desc')
            ->result();
    }

    public function get(): array
    {
        $validate = $this->curd->getValidation();
        if ($validate->fails()) {
            return [
                'error' => 1,
                'msg' => $validate->errors()
            ];
        }

        return $this->curd
            ->getModel('queue')
            ->result();
    }

    /**
     * 新增队列
     * @return array
     */
    public function add(): array
    {
        $body = $this->request->post();
        $validate = $this->curd->addValidation([
            'queue' => 'required',
            'group' => 'required',
            'ttl' => 'required',
            'ttl_unit' => 'required',
            'sizeof' => 'required',
            'sizeof_unit' => 'required',
            'limit' => 'required',
            'limit_unit' => 'required',
            'lazy' => 'required',
            'rewind' => 'required',
            'rewind_ttl' => 'required_if:rewind,true',
            'rewind_ttl_unit' => 'required_if:rewind,true'
        ]);
        if ($validate->fails()) {
            return [
                'error' => 1,
                'msg' => $validate->errors()
            ];
        }
        return $this->curd
            ->addModel('queue', $body)
            ->afterHook(function () use ($body) {
                $master = $this->rabbit->getMaster();
                $masterNode = $this->rabbit->getMasterNode();
                $masterAMQP = $this->rabbit->getMasterAMQP();

                $option = new ExchangeOption();
                $option->setType('topic');
                $response = $master->exchanges()
                    ->put('sys.' . $body['queue'], $option, '/center');
                if ($response->isError()) {
                    Context::set('error', $response->result());
                    Db::rollBack();
                    return false;
                }

                // 创建死信队列
                $option = new QueueOption($masterNode);
                $option->setQueueLazyMode();
                $response = $master->queues()
                    ->put($body['queue'] . '.fails', $option, '/center');
                if ($response->isError()) {
                    Context::set('error', $response->result());
                    Db::rollBack();
                    return false;
                }

                // 绑定系统交换器
                $response = $master->bindings()->setBindingToQueue(
                    'sys.' . $body['queue'],
                    $body['queue'] . '.fails',
                    '/center',
                    'dead'
                );
                if ($response->isError()) {
                    Context::set('error', $response->result());
                    Db::rollBack();
                    return false;
                }

                // 创建队列
                $option = new QueueOption($masterNode);
                $option->setMaxLength($body['limit'] * $body['limit_unit']);
                $option->setDeadLetterExchange('sys.' . $body['queue']);
                $option->setDeadLetterRoutingKey('dead');
                if ($body['ttl'] !== 0) {
                    $option->setMessageTTL($body['ttl'] * $body['ttl_unit'] * 1000);
                }
                if ($body['lazy'] === true) {
                    $option->setQueueLazyMode();
                }
                $response = $master->queues()
                    ->put($body['queue'], $option, '/center');
                if ($response->isError()) {
                    Context::set('error', $response->result());
                    Db::rollBack();
                    return false;
                }

                $broker = $this->rabbit->getBroker('default');
                $brokerNode = $this->rabbit->getBrokerNode('default');
                $brokerAMQP = $this->rabbit->getBrokerAMQP('default');

                // 创建备用分布队列
                $option = new QueueOption($brokerNode);
                $option->setQueueLazyMode();
                $response = $broker->queues()->put($body['queue'] . '.hide', $option, '/center');
                if ($response->isError()) {
                    Context::set('error', $response->result());
                    Db::rollBack();
                    return false;
                }

                // 回溯队列
                if ($body['rewind'] === true) {
                    $option = new QueueOption($masterNode);
                    $option->setMessageTTL($body['rewind_ttl'] * $body['rewind_ttl_unit'] * 1000);
                    $option->setQueueLazyMode();
                    $response = $master->queues()->put($body['queue'] . '.rewind', $option, '/center');
                    if ($response->isError()) {
                        Context::set('error', $response->result());
                        Db::rollBack();
                        return false;
                    }
                    $option = new QueueOption($brokerNode);
                    $option->setMessageTTL(60 * 60 * 24 * 15 * 1000);
                    $option->setQueueLazyMode();
                    $response = $broker->queues()->put($body['queue'] . '.rewind', $option, '/center');
                    if ($response->isError()) {
                        Context::set('error', $response->result());
                        Db::rollBack();
                        return false;
                    }
                }

                // 集群迁移 broker
                $response = $broker->parameters()->put('shovel',
                    $body['queue'],
                    '/center',
                    [
                        'src-delete-after' => 'never',
                        'src-protocol' => 'amqp091',
                        'src-queue' => $body['queue'] . '.fails',
                        'src-uri' => $masterAMQP . '/' . urlencode('/center'),
                        'ack-mode' => 'on-confirm',
                        'dest-add-forward-headers' => false,
                        'dest-protocol' => 'amqp091',
                        'dest-queue' => $body['queue'] . '.hide',
                        'dest-uri' => $brokerAMQP . '/' . urlencode('/center'),
                    ]);

                if ($response->isError()) {
                    Context::set('error', $response->result());
                    Db::rollBack();
                    return false;
                }
                $this->clearRedis();
                return true;
            })
            ->result();
    }

    /**
     * 修改队列
     * @return array
     */
    public function edit(): array
    {
        $body = $this->request->post();
        $validate = $this->curd->editValidation([
            'group' => 'required',
            'ttl' => 'required',
            'ttl_unit' => 'required',
            'sizeof' => 'required',
            'sizeof_unit' => 'required',
            'limit' => 'required',
            'limit_unit' => 'required',
            'lazy' => 'required',
            'rewind' => 'required',
            'rewind_ttl' => 'required_if:rewind,true',
            'rewind_ttl_unit' => 'required_if:rewind,true'
        ]);
        if ($validate->fails()) {
            return [
                'error' => 1,
                'msg' => $validate->errors()
            ];
        }

        $query = Db::table('queue')
            ->where('id', '=', $body['id'])
            ->first();

        return $this->curd
            ->editModel('queue')
            ->afterHook(function () use ($body, $query) {
                $master = $this->rabbit->getMaster();
                $masterNode = $this->rabbit->getMasterNode();

                // 获取原队列信息
                $response = $master->queues()->get($query->queue, '/center');
                if ($response->isError()) {
                    Context::set('error', $response->result());
                    Db::rollBack();
                    return false;
                }
                $data = $response->getData();

                // 获取原队列消息内容
                $option = new MessageOption();
                $option->setCount((int)$data['messages']);
                $response = $master->queues()->getMessage($query->queue, $option, '/center');
                if ($response->isError()) {
                    Context::set('error', $response->result());
                    Db::rollBack();
                    return false;
                }
                $message = $response->getData();

                // 移除原队列
                $response = $master->queues()->delete($query->queue, '/center');
                if ($response->isError()) {
                    Context::set('error', $response->result());
                    Db::rollBack();
                    return false;
                }

                // 创建更新队列
                $option = new QueueOption($masterNode);
                $option->setMaxLength($body['limit'] * $body['limit_unit']);
                $option->setDeadLetterExchange('sys.' . $query->queue);
                $option->setDeadLetterRoutingKey('dead');
                if ($body['ttl'] !== 0) {
                    $option->setMessageTTL($body['ttl'] * $body['ttl_unit'] * 1000);
                }
                if ($body['lazy'] === true) {
                    $option->setQueueLazyMode();
                }
                $response = $master->queues()
                    ->put($query->queue, $option, '/center');
                if ($response->isError()) {
                    Context::set('error', $response->result());
                    Db::rollBack();
                    return false;
                }

                // 恢复队列信息
                foreach ($message as $value) {
                    $option = new PublishOption();
                    $option->setRoutingKey($value['routing_key']);
                    $option->setPayload($value['payload']);
                    $option->setPayloadEncoding($value['payload_encoding']);
                    $option->setProperties($value['properties']);
                    $master->exchanges()->publish('', $option, '/center');
                }

                $broker = $this->rabbit->getBroker('default');
                $brokerNode = $this->rabbit->getBrokerNode('default');
                if ($body['rewind'] === true) {
                    if ($query->rewind === 1 &&
                        ($query->rewind_ttl !== $body['rewind_ttl'] || $query->rewind_ttl_unit !== $body['rewind_ttl_unit'])
                    ) {
                        // 获取原队列信息
                        $response = $master->queues()->get($query->queue . '.rewind', '/center');
                        if ($response->isError()) {
                            Context::set('error', $response->result());
                            Db::rollBack();
                            return false;
                        }
                        $data = $response->getData();
                        // 获取原队列消息内容
                        $option = new MessageOption();
                        $option->setCount((int)$data['messages']);
                        $response = $master->queues()->getMessage($query->queue . '.rewind', $option, '/center');
                        if ($response->isError()) {
                            Context::set('error', $response->result());
                            Db::rollBack();
                            return false;
                        }
                        $message = $response->getData();

                        // 移除原队列
                        $response = $master->queues()->delete($query->queue . '.rewind', '/center');
                        if ($response->isError()) {
                            Context::set('error', $response->result());
                            Db::rollBack();
                            return false;
                        }

                        // 创建更新队列
                        $option = new QueueOption($masterNode);
                        $option->setMessageTTL($body['rewind_ttl'] * $body['rewind_ttl_unit'] * 1000);
                        $option->setQueueLazyMode();
                        $response = $master->queues()
                            ->put($query->queue . '.rewind', $option, '/center');
                        if ($response->isError()) {
                            Context::set('error', $response->result());
                            Db::rollBack();
                            return false;
                        }

                        // 恢复队列信息
                        foreach ($message as $value) {
                            $option = new PublishOption();
                            $option->setRoutingKey($value['routing_key']);
                            $option->setPayload($value['payload']);
                            $option->setPayloadEncoding($value['payload_encoding']);
                            $option->setProperties($value['properties']);
                            $master->exchanges()->publish('', $option, '/center');
                        }

                    }

                    if ($query->rewind === 0) {
                        // 更新回溯队列
                        $option = new QueueOption($masterNode);
                        $option->setMessageTTL($body['rewind_ttl'] * $body['rewind_ttl_unit'] * 1000);
                        $option->setQueueLazyMode();
                        $response = $master->queues()->put($query->queue . '.rewind', $option, '/center');
                        if ($response->isError()) {
                            Context::set('error', $response->result());
                            Db::rollBack();
                            return false;
                        }
                        $option = new QueueOption($brokerNode);
                        $option->setMessageTTL(60 * 60 * 24 * 15 * 1000);
                        $response = $broker->queues()->put($query->queue . '.rewind', $option, '/center');
                        if ($response->isError()) {
                            Context::set('error', $response->result());
                            Db::rollBack();
                            return false;
                        }
                    }
                } else {
                    $master->queues()
                        ->delete($query->queue . '.rewind', '/center');
                    $broker->queues()
                        ->delete($query->queue . '.rewind', '/center');
                }
                // 绑定的重置恢复
                $bindings = Db::table('topic_binding')
                    ->where('queue', '=', $query->queue)
                    ->get();
                if (!empty($bindings)) {
                    foreach ($bindings as $value) {
                        $response = $master->bindings()->setBindingToQueue(
                            $value->topic,
                            $value->queue,
                            '/center',
                            ''
                        );
                        if ($response->isError()) {
                            Context::set('error', $response->result());
                            Db::rollBack();
                            return false;
                        }
                    }
                }
                // 代理订阅器的重置恢复
                $agent = Db::table('queue_agent')
                    ->where('queue', '=', $query->queue)
                    ->first();
                if (!empty($agent)) {
                    $response = $this->amqpSubscriberService->put(
                        (string)$agent->id,
                        $agent->queue,
                        $agent->url,
                        $agent->secret
                    );
                    if ($response->getError() !== 0) {
                        Context::set('error', [
                            'error' => 1,
                            'msg' => "<{$agent->id}> " . $response->getMsg()
                        ]);
                        Db::rollBack();
                        return false;
                    }
                }
                $this->clearRedis();
                return true;
            })
            ->result();
    }

    /**
     * 删除队列
     * @return array
     */
    public function delete(): array
    {
        $body = $this->request->post();
        $validate = $this->curd->deleteValidation();
        if ($validate->fails()) {
            return [
                'error' => 1,
                'msg' => $validate->errors()
            ];
        }

        $query = Db::table('queue')
            ->whereIn('id', $body['id'])
            ->get();

        return $this->curd
            ->deleteModel('queue')
            ->afterHook(function () use ($query) {
                $master = $this->rabbit->getMaster();
                $broker = $this->rabbit->getBroker('default');
                foreach ($query->toArray() as $value) {
                    $master
                        ->queues()
                        ->delete($value->queue . '.fails', '/center');
                    $master
                        ->exchanges()
                        ->delete('sys.' . $value->queue, '/center');
                    $master
                        ->queues()
                        ->delete($value->queue, '/center');
                    $master
                        ->queues()
                        ->delete($value->queue . '.rewind', '/center');
                    $broker
                        ->queues()
                        ->delete($value->queue . '.hide', '/center');
                    $broker
                        ->queues()
                        ->delete($value->queue . '.rewind', '/center');
                    $broker
                        ->parameters()
                        ->delete(
                            'shovel',
                            '/center',
                            $value->queue
                        );
                }
                $this->clearRedis();
                return true;
            })
            ->result();
    }

    /**
     * 清除缓存
     */
    private function clearRedis(): void
    {
        $this->queueRedis->clear();
    }

    /**
     * 获取队列信息
     * @return array
     * @throws Exception
     */
    public function getQueueInfo(): array
    {
        $body = $this->request->post();
        $validator = $this->validation->make($body, [
            'queue' => 'required|array',
        ]);
        if ($validator->fails()) {
            return [
                'error' => 1,
                'msg' => $validator->errors()
            ];
        }
        $master = $this->rabbit->getMaster();
        $response = $master->queues()->lists('/center');
        if ($response->isError()) {
            return $response->result();
        }
        $queue = $body['queue'];
        $queueFails = array_map(fn($v) => $v . '.fails', $queue);
        $lists = [];
        foreach ($response->getData() as $value) {
            if (in_array($value['name'], $queue, true)) {
                $lists[$value['name']] = [
                    'ready' => $value['messages_ready'] ?? 0,
                    'ready_rate' => $value['messages_ready_details']['rate'] ?? 0,
                    'unacknowledged' => $value['messages_unacknowledged'] ?? 0,
                    'unacknowledged_rate' => $value['messages_unacknowledged_details']['rate'] ?? 0
                ];
            }
            if (in_array($value['name'], $queueFails, true)) {
                $name = str_replace('.fails', '', $value['name']);
                $lists[$name]['dead'] = $value['messages_ready'];
                $lists[$name]['dead_rate'] = $value['messages_ready_details']['rate'];
            }
        }
        $queueHide = array_map(fn($v) => $v . '.hide', $queue);
        $broker = $this->rabbit->getBroker('default');
        $response = $broker->queues()->lists('/center');
        if ($response->isError()) {
            return $response->result();
        }
        foreach ($response->getData() as $value) {
            if (in_array($value['name'], $queueHide, true)) {
                $name = str_replace('.hide', '', $value['name']);
                $lists[$name]['dead'] += $value['messages_ready'];
                $lists[$name]['dead_rate'] += $value['messages_ready_details']['rate'];
            }
        }
        return [
            'error' => 0,
            'data' => $lists
        ];
    }

    /**
     * 获取队列消息内容
     * @return array
     * @throws Exception
     */
    public function receive(): array
    {
        $body = $this->request->post();
        $validator = $this->validation->make($body, [
            'queue' => 'required',
            'count' => 'required|numeric'
        ]);
        if ($validator->fails()) {
            return [
                'error' => 1,
                'msg' => $validator->errors()
            ];
        }
        $master = $this->rabbit->getMaster();
        $option = new MessageOption();
        $option->setCount($body['count']);
        return $master
            ->queues()
            ->getMessage($body['queue'], $option, '/center')
            ->result();
    }

    /**
     * 获取回溯队列
     * @return array
     * @throws Exception
     */
    public function getRewind(): array
    {
        $body = $this->request->post();
        $validator = $this->validation->make($body, [
            'queue' => 'required',
        ]);
        if ($validator->fails()) {
            return [
                'error' => 1,
                'msg' => $validator->errors()
            ];
        }
        $master = $this->rabbit->getMaster();
        $response = $master->queues()->get($body['queue'] . '.rewind', '/center');
        if ($response->isError()) {
            return $response->result();
        }
        $count = (int)$response->getData()['messages'];
        if ($count === 0) {
            return [
                'error' => 0,
                'data' => []
            ];
        }
        $option = new MessageOption();
        $option->setCount($count);
        return $master
            ->queues()
            ->getMessage($body['queue'] . '.rewind', $option, '/center')
            ->result();
    }

    /**
     * 执行回溯处理
     * @return array
     * @throws Exception
     */
    public function startRewind(): array
    {
        $body = $this->request->post();
        $validator = $this->validation->make($body, [
            'queue' => 'required',
        ]);
        if ($validator->fails()) {
            return [
                'error' => 1,
                'msg' => $validator->errors()
            ];
        }
        $master = $this->rabbit->getMaster();
        $response = $master->queues()->get($body['queue'] . '.rewind', '/center');
        if ($response->isError()) {
            return $response->result();
        }
        $count = (int)$response->getData()['messages'];
        if ($count === 0) {
            return [
                'error' => 0,
                'msg' => 'ok'
            ];
        }
        $option = new MessageOption();
        $option->setCount($count);
        $option->setAck(true);
        $response = $master->queues()
            ->getMessage($body['queue'] . '.rewind', $option, '/center');
        if ($response->isError()) {
            return $response->result();
        }
        foreach ($response->getData() as $value) {
            $option = new PublishOption();
            $option->setRoutingKey($body['queue']);
            $option->setPayload($value['payload']);
            $option->setPayloadEncoding($value['payload_encoding']);
            $option->setProperties($value['properties']);
            $response = $master->exchanges()->publish('', $option, '/center');
            if ($response->isError()) {
                return $response->result();
            }
        }
        return [
            'error' => 0,
            'msg' => 'ok'
        ];
    }

    /**
     * 日志搜索
     * @return array
     */
    public function logs(): array
    {
        $body = $this->request->post();
        $validator = $this->validation->make($body, [
            'page' => 'required',
            'page.limit' => 'required|integer|between:1,50',
            'page.index' => 'required|integer|min:1',
            'queue' => 'required|string',
            'get' => 'sometimes|array|size:2'
        ]);
        if ($validator->fails()) {
            return [
                'error' => 1,
                'msg' => $validator->errors()
            ];
        }
        $response = $this->es->client()->search([
            'index' => $this->config->get('elasticsearch.index.mq-message'),
            'body' => [
                'query' => empty($body['get']) ? [
                    'match' => [
                        'queue' => $body['queue']
                    ]
                ] : [
                    'bool' => [
                        'must' => [
                            [
                                'match' => [
                                    'queue' => $body['queue']
                                ]
                            ],
                            [
                                'range' => [
                                    'get' => [
                                        'gte' => $body['get'][0],
                                        'lte' => $body['get'][1]
                                    ]
                                ]
                            ]
                        ]
                    ]
                ],
                'sort' => [
                    'get' => 'desc'
                ]
            ],
            'from' => ($body['page']['index'] - 1) * $body['page']['limit'],
            'size' => $body['page']['limit'],
        ]);
        return [
            'error' => 0,
            'data' => [
                'lists' => $response['hits']['hits'] ?? [],
                'total' => $response['hits']['total']['value'] ?? 0,
            ]
        ];
    }

    /**
     * Exists Queue
     * @return array
     */
    public function validedQueue(): array
    {
        $body = $this->request->post();
        if (empty($body['queue'])) {
            return [
                'error' => 1,
                'msg' => 'require queue'
            ];
        }

        $exists = Db::table('queue')
            ->where('queue', '=', $body['queue'])
            ->exists();

        return [
            'error' => 0,
            'data' => $exists
        ];
    }

}