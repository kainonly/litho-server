<?php
declare(strict_types=1);

namespace App\Controller\System;

use App\Client\ElasticSearchClient;
use App\HttpClient\RabbitMQService;
use App\RedisModel\Common\TopicRedis;
use Exception;
use Hyperf\DbConnection\Db;
use Hyperf\Di\Annotation\Inject;
use Hyperf\Utils\Context;
use RabbitMQ\API\Common\ExchangeOption;
use RabbitMQ\API\Common\PublishOption;

class TopicController extends BaseController
{
    /**
     * @Inject()
     * @var RabbitMQService
     */
    private RabbitMQService $rabbit;
    /**
     * @Inject()
     * @var TopicRedis
     */
    private TopicRedis $topicRedis;
    /**
     * @Inject()
     * @var ElasticSearchClient
     */
    private ElasticSearchClient $es;

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
            ->originListsModel('topic')
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
            ->listsModel('topic')
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
            ->getModel('topic')
            ->result();
    }

    public function add(): array
    {
        $body = $this->request->post();
        $validate = $this->curd->addValidation([
            'topic' => 'required',
            'group' => 'required',
            'sizeof' => 'required',
            'sizeof_unit' => 'required'
        ]);
        if ($validate->fails()) {
            return [
                'error' => 1,
                'msg' => $validate->errors()
            ];
        }
        return $this->curd
            ->addModel('topic', $body)
            ->afterHook(function () use ($body) {
                $master = $this->rabbit->getMaster();
                $option = new ExchangeOption();
                $option->setType('topic');
                $response = $master->exchanges()->put(
                    $body['topic'],
                    $option,
                    '/center'
                );
                if ($response->isError()) {
                    Context::set('error', $response->result());
                    return false;
                }
                $this->clearRedis();
                return true;
            })
            ->result();
    }

    public function edit(): array
    {
        $validate = $this->curd->editValidation([
            'group' => 'required',
            'sizeof' => 'required',
            'sizeof_unit' => 'required'
        ]);
        if ($validate->fails()) {
            return [
                'error' => 1,
                'msg' => $validate->errors()
            ];
        }
        return $this->curd
            ->editModel('topic')
            ->afterHook(function () {
                $this->clearRedis();
                return true;
            })
            ->result();
    }

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

        $query = Db::table('topic')
            ->whereIn('id', $body['id'])
            ->get();

        return $this->curd
            ->deleteModel('topic')
            ->afterHook(function () use ($query) {
                $master = $this->rabbit->getMaster();
                foreach ($query->toArray() as $value) {
                    $master->exchanges()->delete($value->topic, '/center');
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
        $this->topicRedis->clear();
    }

    /**
     * 主题绑定获取
     * @return array
     */
    public function getBinding(): array
    {
        $body = $this->request->post();
        $validator = $this->validation->make($body, [
            'topic' => 'required',
        ]);
        if ($validator->fails()) {
            return [
                'error' => 1,
                'msg' => $validator->errors()
            ];
        }

        $query = Db::table('topic_binding')
            ->where('topic', '=', $body['topic'])
            ->get();

        return [
            'error' => 0,
            'data' => $query->toArray()
        ];
    }

    /**
     * 主题绑定更新
     * @return array
     */
    public function putBinding(): array
    {
        $body = $this->request->post();
        $validator = $this->validation->make($body, [
            'topic' => 'required',
            'queue' => 'required|array',
        ]);
        if ($validator->fails()) {
            return [
                'error' => 1,
                'msg' => $validator->errors()
            ];
        }

        return !Db::transaction(function () use ($body) {
            Db::table('topic_binding')
                ->insert(
                    array_map(fn($v) => [
                        'topic' => $body['topic'],
                        'queue' => $v
                    ], $body['queue'])
                );
            $master = $this->rabbit->getMaster();
            foreach ($body['queue'] as $value) {
                $response = $master->bindings()->setBindingToQueue(
                    $body['topic'],
                    $value,
                    '/center',
                    ''
                );
                if ($response->isError()) {
                    Db::rollBack();
                    return false;
                }
            }
            return true;
        }) ? [
            'error' => 1,
            'msg' => 'topic binding failed'
        ] : [
            'error' => 0,
            'msg' => 'ok'
        ];
    }

    /**
     * 主题绑定删除
     * @return array
     */
    public function deleteBinding(): array
    {
        $body = $this->request->post();
        $validator = $this->validation->make($body, [
            'id' => 'required',
        ]);
        if ($validator->fails()) {
            return [
                'error' => 1,
                'msg' => $validator->errors()
            ];
        }

        $query = Db::table('topic_binding')
            ->where('id', '=', $body['id'])
            ->first();

        return !Db::transaction(function () use ($body, $query) {
            Db::table('topic_binding')->delete($body['id']);
            $master = $this->rabbit->getMaster();
            $response = $master->bindings()->deleteBindingToQueueFormRoutingKey(
                $query->topic,
                $query->queue,
                '/center'
            );

            if ($response->isError()) {
                Db::rollBack();
                return false;
            }
            return true;
        }) ? [
            'error' => 1,
            'msg' => 'Binding delete failed'
        ] : [
            'error' => 0,
            'msg' => 'ok'
        ];
    }

    /**
     * 获取主题信息
     * @return array
     * @throws Exception
     */
    public function getTopicInfo(): array
    {
        $body = $this->request->post();
        $validator = $this->validation->make($body, [
            'topic' => 'required|array',
        ]);
        if ($validator->fails()) {
            return [
                'error' => 1,
                'msg' => $validator->errors()
            ];
        }
        $master = $this->rabbit->getMaster();
        $response = $master->exchanges()->lists('/center');
        if ($response->isError()) {
            return $response->result();
        }
        $lists = [];
        foreach ($response->getData() as $value) {
            if (in_array($value['name'], $body['topic'], true)) {
                if (!empty($value['message_stats'])) {
                    $lists[$value['name']] = [
                        'in' => $value['message_stats']['publish_in'] ?? 0,
                        'in_rate' => $value['message_stats']['publish_in_details']['rate'] ?? 0,
                        'out' => $value['message_stats']['publish_out'] ?? 0,
                        'out_rate' => $value['message_stats']['publish_out_details']['rate'] ?? 0
                    ];
                } else {
                    $lists[$value['name']] = [
                        'in' => 0,
                        'in_rate' => 0,
                        'out' => 0,
                        'out_rate' => 0
                    ];
                }
            }
        }
        return [
            'error' => 0,
            'data' => $lists
        ];
    }

    /**
     * 投递消息
     * @return array
     * @throws Exception
     */
    public function publish(): array
    {
        $body = $this->request->post();
        $validator = $this->validation->make($body, [
            'topic' => 'required',
            'payload' => 'required'
        ]);
        if ($validator->fails()) {
            return [
                'error' => 1,
                'msg' => $validator->errors()
            ];
        }
        $master = $this->rabbit->getMaster();
        $option = new PublishOption();
        $option->setPayload($body['payload']);
        $option->setRoutingKey('');
        $response = $master->exchanges()
            ->publish($body['topic'], $option, '/center');
        return $response->result();
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
            'topic' => 'required|string',
            'time' => 'sometimes|array|size:2'
        ]);
        if ($validator->fails()) {
            return [
                'error' => 1,
                'msg' => $validator->errors()
            ];
        }
        $response = $this->es->client()->search([
            'index' => $this->config->get('elasticsearch.index.mq-publish'),
            'body' => [
                'query' => empty($body['time']) ? [
                    'match' => [
                        'topic' => $body['topic']
                    ]
                ] : [
                    'bool' => [
                        'must' => [
                            [
                                'match' => [
                                    'topic' => $body['topic']
                                ]
                            ],
                            [
                                'range' => [
                                    'time' => [
                                        'gte' => $body['time'][0],
                                        'lte' => $body['time'][1]
                                    ]
                                ]
                            ]
                        ]
                    ]
                ],
                'sort' => [
                    'time' => 'desc'
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
     * Exists Topic
     * @return array
     */
    public function validedTopic(): array
    {
        $body = $this->request->post();
        if (empty($body['topic'])) {
            return [
                'error' => 1,
                'msg' => 'require topic'
            ];
        }

        $exists = Db::table('topic')
            ->where('topic', '=', $body['topic'])
            ->exists();

        return [
            'error' => 0,
            'data' => $exists
        ];
    }
}