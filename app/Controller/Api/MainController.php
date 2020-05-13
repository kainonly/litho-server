<?php
declare(strict_types=1);

namespace App\Controller\Api;

class MainController extends BaseController
{
    public function alivetest(): array
    {
        $body = $this->request->post();
        $validator = $this->validation->make($body, [
            'uuid' => 'required|uuid'
        ]);
        if ($validator->fails()) {
            return [
                'error' => 1,
                'msg' => $validator->errors()
            ];
        }
        return [
            'error' => 0,
            'data' => [
                'uuid' => $body['uuid']
            ]
        ];
    }

}