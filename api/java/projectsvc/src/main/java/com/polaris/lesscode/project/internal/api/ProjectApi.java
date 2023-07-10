package com.polaris.lesscode.project.internal.api;


import com.polaris.lesscode.project.internal.req.AppEvent;
import com.polaris.lesscode.project.internal.req.ChangeProjectChatMemberReq;
import com.polaris.lesscode.project.internal.resp.ApplyProjectTemplateResp;
import com.polaris.lesscode.project.internal.resp.ProjectTemplate;
import com.polaris.lesscode.vo.Result;
import org.springframework.web.bind.annotation.*;

import java.util.List;


@RequestMapping("/api/projectsvc")
public interface ProjectApi {

    @GetMapping(value = "/getProjectTemplateInner")
    Result<ProjectTemplate> getProjectTemplateInner(@RequestParam(value="orgId") Long orgId, @RequestParam("projectId") Long projectId);

    @PostMapping(value = "/applyProjectTemplateInner")
    Result<ApplyProjectTemplateResp> applyProjectTemplateInner(@RequestParam(value="orgId") Long orgId, @RequestParam("userId") Long userId, @RequestParam("appId") Long appId, @RequestBody ProjectTemplate projectTemplate);

    @PostMapping(value = "/changeProjectChatMember")
    Result<?> changeProjectChatMember(@RequestParam(value="orgId") Long orgId, @RequestParam("userId") Long userId, @RequestBody ChangeProjectChatMemberReq changeProjectChatMemberReq);

    @PostMapping(value = "/deleteProjectBatchInner")
    Result<?> deleteProjectBatchInner(@RequestParam(value="orgId") Long orgId, @RequestParam("userId") Long userId, @RequestBody List<Long> projectIds);

    @GetMapping(value = "/authCreateProject")
    Result<?> authCreateProject(@RequestParam(value="orgId") Long orgId);

    @PostMapping(value = "/reportAppEvent")
    Result<?> reportAppEvent(@RequestParam(value="eventType") Integer eventType, @RequestParam("traceId") String traceId, @RequestBody AppEvent req);
}
