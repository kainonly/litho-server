<?php
declare(strict_types=1);
return [
    'schedule' => [
        'host' => env('MICRO_SCHEDULE_HOST'),
        'max_connections' => 50
    ],
    'ssh' => [
        'host' => env('MICRO_SSH_HOST'),
        'max_connections' => 100
    ],
    'amqp' => [
        'host' => env('MICRO_AMQP_HOST'),
        'max_connections' => 10
    ]
];