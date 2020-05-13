<?php
declare(strict_types=1);

namespace App\Controller\Api;

use Hyperf\Di\Annotation\Inject;
use Hyperf\HttpServer\Contract\RequestInterface;
use Hyperf\HttpServer\Contract\ResponseInterface;
use Hyperf\Validation\Contract\ValidatorFactoryInterface;

class BaseController
{
    /**
     * @Inject()
     * @var RequestInterface
     */
    protected RequestInterface $request;
    /**
     * @Inject()
     * @var ResponseInterface
     */
    protected ResponseInterface $response;
    /**
     * @Inject()
     * @var ValidatorFactoryInterface
     */
    protected ValidatorFactoryInterface $validation;
}