/**
 *
 */
package com.polaris.lesscode.project.internal.fallback;

import com.polaris.lesscode.feign.AbstractBaseFallback;
import com.polaris.lesscode.project.internal.api.ProjectApi;
import com.polaris.lesscode.project.internal.req.AppEvent;
import com.polaris.lesscode.project.internal.req.ChangeProjectChatMemberReq;
import com.polaris.lesscode.project.internal.resp.ApplyProjectTemplateResp;
import com.polaris.lesscode.project.internal.resp.ProjectTemplate;
import com.polaris.lesscode.vo.Result;
import feign.hystrix.FallbackFactory;
import org.springframework.stereotype.Component;

import java.util.List;


/**
 * @author admin
 *
 */
@Component
public class ProjectFallbackFactory extends AbstractBaseFallback implements FallbackFactory<ProjectApi> {

    @Override
    public ProjectApi create(Throwable cause) {
        return new ProjectApi() {
            @Override
            public Result<ProjectTemplate> getProjectTemplateInner(Long orgId, Long projectId) {
                cause.printStackTrace();
                return Result.ok(new ProjectTemplate());
            }

            @Override
            public Result<ApplyProjectTemplateResp> applyProjectTemplateInner(Long orgId, Long userId, Long appId, ProjectTemplate projectTemplate) {
                return Result.ok(new ApplyProjectTemplateResp());
            }

            @Override
            public Result<?> changeProjectChatMember(Long orgId, Long userId, ChangeProjectChatMemberReq changeProjectChatMemberReq) {
                return Result.ok();
            }

            @Override
            public Result<?> deleteProjectBatchInner(Long orgId, Long userId, List<Long> projectIds) {
                return Result.ok();
            }

            @Override
            public Result<?> authCreateProject(Long orgId) {
                return Result.ok();
            }

            @Override
            public Result<?> reportAppEvent(Integer eventType, String traceId, AppEvent appEvent) {
                return Result.ok();
            }
        };
    }

}
