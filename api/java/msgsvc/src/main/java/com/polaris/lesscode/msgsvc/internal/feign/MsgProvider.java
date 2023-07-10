/**
 * 
 */
package com.polaris.lesscode.msgsvc.internal.feign;

import com.polaris.lesscode.msgsvc.internal.api.MsgApi;
import org.springframework.cloud.openfeign.FeignClient;

@FeignClient(name = "msgsvc")
public interface MsgProvider extends MsgApi {

}
