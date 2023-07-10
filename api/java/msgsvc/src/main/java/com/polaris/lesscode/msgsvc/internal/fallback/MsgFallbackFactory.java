/**
 *
 */
package com.polaris.lesscode.msgsvc.internal.fallback;

import com.polaris.lesscode.feign.AbstractBaseFallback;
import com.polaris.lesscode.msgsvc.internal.api.MsgApi;
import com.polaris.lesscode.msgsvc.internal.req.PushMsgToMqttReq;
import com.polaris.lesscode.vo.Result;
import feign.hystrix.FallbackFactory;
import org.springframework.stereotype.Component;

import java.util.List;


/**
 * @author admin
 *
 */
@Component
public class MsgFallbackFactory extends AbstractBaseFallback implements FallbackFactory<MsgApi> {

    @Override
    public MsgApi create(Throwable cause) {
        return new MsgApi() {

            @Override
            public Result<?> pushMsgToMqtt(PushMsgToMqttReq req) {
                return null;
            }
        };
    }

}
