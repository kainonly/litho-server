<?php
declare(strict_types=1);

namespace App\Controller\System;

use AMQPSubscriber\Response;
use App\Client\ElasticSearchClient;
use App\GrpcClient\AMQPSubscriberServiceInterface;
use Hyperf\Curd\Common\AddAfterParams;
use Hyperf\Curd\Common\DeleteAfterParams;
use Hyperf\Curd\Common\EditAfterParams;
use Hyperf\DbConnection\Db;
use Hyperf\Di\Annotation\Inject;
use Hyperf\Utils\Context;

class QueueAgentController extends BaseController
{
    /**
     * @Inject()
     * @var AMQPSubscriberServiceInterface
     */
    private AMQPSubscriberServiceInterface $amqpSubscriberService;
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
            ->originListsModel('queue_agent')
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
            ->listsModel('queue_agent')
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
            ->getModel('queue_agent')
            ->result();
    }

    public function add(): array
    {
        $body = $this->request->post();
        $validate = $this->curd->addValidation([
            'name' => 'required',
            'queue' => 'required',
            'url' => 'required',
        ]);
        if ($validate->fails()) {
            return [
                'error' => 1,
                'msg' => $validate->errors()
            ];
        }

        return $this->curd
            ->addModel('queue_agent', $body)
            ->afterHook(function (AddAfterParams $params) use ($body) {
                $response = $this->afterRequest($params->getId(), $body);
                $this->amqpSubscriberService->close();
                if ($response->getError() !== 0) {
                    Context::set('error', [
                        'error' => 1,
                        'msg' => $response->getMsg()
                    ]);
                    Db::rollBack();
                    return false;
                }
                return true;
            })
            ->result();
    }

    public function edit(): array
    {
        $body = $this->request->post();
        $validate = $this->curd->editValidation([
            'name' => 'required_if:switch,false',
            'queue' => 'required_if:switch,false',
            'url' => 'required_if:switch,false',
        ]);
        if ($validate->fails()) {
            return [
                'error' => 1,
                'msg' => $validate->errors()
            ];
        }

        return $this->curd
            ->editModel('queue_agent', $body)
            ->afterHook(function (EditAfterParams $params) use ($body) {
                if (!$params->isSwitch()) {
                    $response = $this->afterRequest($params->getId(), $body);
                    $this->amqpSubscriberService->close();
                    if ($response->getError() !== 0) {
                        Context::set('error', [
                            'error' => 1,
                            'msg' => $response->getMsg()
                        ]);
                        return false;
                    }
                }
                return true;
            })
            ->result();
    }

    public function delete(): array
    {
        $validate = $this->curd->deleteValidation();
        if ($validate->fails()) {
            return [
                'error' => 1,
                'msg' => $validate->errors()
            ];
        }

        return $this->curd
            ->deleteModel('queue_agent')
            ->afterHook(function (DeleteAfterParams $params) {
                $response = $this->amqpSubscriberService->delete((string)$params->getId()[0]);
                $this->amqpSubscriberService->close();
                if ($response->getError() !== 0) {
                    Context::set('error', [
                        'error' => 1,
                        'msg' => $response->getMsg()
                    ]);
                    return false;
                }
                return true;
            })
            ->result();
    }

    /**
     * @param int $id
     * @param array $body
     * @return Response
     */
    private function afterRequest(int $id, array $body): Response
    {
        return $this->amqpSubscriberService->put(
            (string)$id,
            $body['queue'],
            $body['url'],
            $body['secret']
        );
    }

    /**
     * 获取正在运行的代理器
     * @return array
     */
    public function runtime(): array
    {
        $response = $this->amqpSubscriberService->all();
        $lists = [];
        foreach ($response->getData()->getIterator() as $value) {
            $lists[] = $value;
        }
        $this->amqpSubscriberService->close();
        return [
            'error' => 0,
            'data' => $lists
        ];
    }

    /**
     * 服务同步
     * @return array
     */
    public function sync(): array
    {
        $response = $this->amqpSubscriberService->all();
        if ($response->getError() !== 0) {
            return [
                'error' => 1,
                'msg' => $response->getMsg()
            ];
        }
        foreach ($response->getData()->getIterator() as $identity) {
            $deleteResponse = $this->amqpSubscriberService->delete((string)$identity);
            if ($deleteResponse->getError() !== 0) {
                return [
                    'error' => 1,
                    'msg' => "<$identity> " . $deleteResponse->getMsg()
                ];
            }
        }
        $query = Db::table('queue_agent')->get();
        foreach ($query->toArray() as $value) {
            $putResponse = $this->amqpSubscriberService->put(
                (string)$value->id,
                $value->queue,
                $value->url,
                $value->secret
            );
            if ($putResponse->getError() !== 0) {
                return [
                    'error' => 1,
                    'msg' => "<{$value->id}> " . $putResponse->getMsg()
                ];
            }
        }
        $this->amqpSubscriberService->close();
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
            'identity' => 'required|string',
            'time' => 'sometimes|array|size:2'
        ]);
        if ($validator->fails()) {
            return [
                'error' => 1,
                'msg' => $validator->errors()
            ];
        }
        $response = $this->es->client()->search([
            'index' => $this->config->get('elasticsearch.index.amqp-subscriber'),
            'body' => [
                'query' => empty($body['time']) ? [
                    'match' => [
                        'Identity' => $body['identity']
                    ]
                ] : [
                    'bool' => [
                        'must' => [
                            [
                                'match' => [
                                    'Identity' => $body['identity']
                                ]
                            ],
                            [
                                'range' => [
                                    'Time' => [
                                        'gte' => $body['time'][0],
                                        'lte' => $body['time'][1]
                                    ]
                                ]
                            ]
                        ]
                    ]
                ],
                'sort' => [
                    'Time' => 'desc'
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
}
