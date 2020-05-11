<?php
declare(strict_types=1);

namespace App\Controller\System;

use App\Client\ElasticSearchClient;
use App\GrpcClient\ScheduleServiceInterface;
use Hyperf\Curd\Common\AddAfterParams;
use Hyperf\Curd\Common\DeleteAfterParams;
use Hyperf\Curd\Common\EditAfterParams;
use Hyperf\DbConnection\Db;
use Hyperf\Di\Annotation\Inject;
use Hyperf\Utils\Context;
use ScheduleMicroservice\EntryOptionWithTime;
use ScheduleMicroservice\Information;
use ScheduleMicroservice\Response;

class ScheduleController extends BaseController
{
    /**
     * @Inject()
     * @var ScheduleServiceInterface
     */
    private ScheduleServiceInterface $scheduleService;
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
            ->originListsModel('schedule')
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
            ->listsModel('schedule')
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
            ->getModel('schedule')
            ->result();
    }

    public function add(): array
    {
        $body = $this->request->post();
        $validate = $this->curd->addValidation([
            'name' => 'required',
            'cron_time' => 'required',
            'zone' => 'required',
            'group' => 'required',
            'url' => 'required'
        ]);
        if ($validate->fails()) {
            return [
                'error' => 1,
                'msg' => $validate->errors()
            ];
        }
        return $this->curd
            ->addModel('schedule')
            ->afterHook(function (AddAfterParams $params) use ($body) {
                $response = $this->afterRequest((string)$params->getId(), $body);
                $this->scheduleService->close();
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

    public function edit(): array
    {
        $body = $this->request->post();
        $validate = $this->curd->editValidation([
            'name' => 'required_if:switch,false',
            'cron_time' => 'required_if:switch,false',
            'zone' => 'required_if:switch,false',
            'group' => 'required_if:switch,false',
            'url' => 'required_if:switch,false'
        ]);
        if ($validate->fails()) {
            return [
                'error' => 1,
                'msg' => $validate->errors()
            ];
        }
        return $this->curd
            ->editModel('schedule')
            ->afterHook(function (EditAfterParams $params) use ($body) {
                if (!$params->isSwitch()) {
                    $response = $this->afterRequest((string)$params->getId(), $body);
                } else {
                    $response = $this->scheduleService->running(
                        (string)$params->getId(),
                        (bool)$body['status']
                    );
                }
                $this->scheduleService->close();
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
            ->deleteModel('schedule')
            ->afterHook(function (DeleteAfterParams $params) {
                $response = $this->scheduleService->delete((string)$params->getId()[0]);
                $this->scheduleService->close();
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
     * @param string $id
     * @param array $body
     * @return Response
     */
    private function afterRequest(string $id, array $body): Response
    {
        return $this->scheduleService->put(
            $id,
            $body['zone'],
            (bool)$body['status'],
            [
                $id => [
                    'cron_time' => $body['cron_time'],
                    'url' => $body['url'],
                    'headers' => [
                        'secret' => $body['secret']
                    ],
                    'body' => (object)[],
                ]
            ]
        );
    }

    /**
     * 获取任务调度器信息
     * @return array
     */
    public function listsJobs(): array
    {
        $body = $this->request->post();
        $validator = $this->validation->make($body, [
            'identity' => 'required',
        ]);
        if ($validator->fails()) {
            return [
                'error' => 1,
                'msg' => $validator->errors()
            ];
        }
        $response = $this->scheduleService->lists($body['identity']);
        if ($response->getError() !== 0) {
            return [
                'error' => 1,
                'msg' => $response->getMsg()
            ];
        }
        $lists = [];
        /**
         * @var Information $value
         */
        foreach ($response->getData()->getIterator() as $value) {
            $entries = [];
            /**
             * @var EntryOptionWithTime $entry
             */
            foreach ($value->getEntries()->getIterator() as $task => $entry) {
                $entries[$task] = [
                    'cron_time' => $entry->getCronTime(),
                    'url' => $entry->getUrl(),
                    'headers' => json_decode($entry->getHeaders()),
                    'body' => json_decode($entry->getBody()),
                    'next_date' => $entry->getNextDate(),
                    'last_date' => $entry->getLastDate()
                ];
            }
            $lists[] = [
                'identity' => (int)$value->getIdentity(),
                'time_zone' => $value->getTimeZone(),
                'start' => $value->getStart(),
                'entries' => $entries
            ];
        }
        $this->scheduleService->close();
        return [
            'error' => 0,
            'data' => $lists
        ];
    }

    /**
     * 同步服务
     * @return array
     */
    public function sync(): array
    {
        $response = $this->scheduleService->all();
        if ($response->getError() !== 0) {
            return [
                'error' => 1,
                'msg' => $response->getMsg()
            ];
        }
        foreach ($response->getData()->getIterator() as $identity) {
            $deleteResponse = $this->scheduleService->delete((string)$identity);
            if ($deleteResponse->getError() !== 0) {
                return [
                    'error' => 1,
                    'msg' => "<$identity> " . $deleteResponse->getMsg()
                ];
            }
        }
        $query = Db::table('schedule')->get();
        foreach ($query->toArray() as $value) {
            $putResponse = $this->scheduleService->put(
                (string)$value->id,
                $value->zone,
                true,
                [
                    (string)$value->id => [
                        'cron_time' => $value->cron_time,
                        'url' => $value->url,
                        'headers' => [
                            'secret' => $value->secret
                        ],
                        'body' => [],
                    ]
                ]
            );
            if ($putResponse->getError() !== 0) {
                return [
                    'error' => 1,
                    'msg' => "<{$value->id}> " . $putResponse->getMsg()
                ];
            }
        }
        $this->scheduleService->close();
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
            'index' => $this->config->get('elasticsearch.index.schedule'),
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