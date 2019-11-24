package van.api.controllers

import com.alibaba.fastjson.JSONObject
import org.springframework.boot.autoconfigure.EnableAutoConfiguration
import org.springframework.web.bind.annotation.GetMapping
import org.springframework.web.bind.annotation.RestController

@RestController
@EnableAutoConfiguration
class IndexController {
    @GetMapping("/")
    fun index(): String {
        val context = JSONObject()
        context["project"] = "van-api"
        context["author"] = "kain"
        return context.toJSONString()
    }
}