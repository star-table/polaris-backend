/**
 * 
 */
package com.polaris.lesscode.project.internal.feign;

import org.springframework.cloud.openfeign.FeignClient;

import com.polaris.lesscode.project.internal.api.ProjectApi;
import com.polaris.lesscode.project.internal.fallback.ProjectFallbackFactory;

@FeignClient(name = "projectsvc")
public interface ProjectProvider extends ProjectApi {

}
