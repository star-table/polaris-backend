package consts

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

//Issue 关联类型
//负责人
func TestConstIssueRelationTypeOwner(t *testing.T) {

	convey.Convey("Test IssueRelationTypeOwner", t, func() {
		convey.Convey("IssueRelationTypeOwner", func() {
			convey.So(IssueRelationTypeOwner, convey.ShouldEqual, 1)
		})
	})
}

//参与人
func TestConstIssueRelationTypeParticipant(t *testing.T) {

	convey.Convey("Test IssueRelationTypeParticipant", t, func() {
		convey.Convey("IssueRelationTypeParticipant", func() {
			convey.So(IssueRelationTypeParticipant, convey.ShouldEqual, 2)
		})
	})
}

//关注人
func TestConstIssueRelationTypeFollower(t *testing.T) {

	convey.Convey("Test IssueRelationTypeFollower", t, func() {
		convey.Convey("IssueRelationTypeFollower", func() {
			convey.So(IssueRelationTypeFollower, convey.ShouldEqual, 3)
		})
	})
}

//关联任务
func TestConstIssueRelationTypeIssue(t *testing.T) {

	convey.Convey("Test IssueRelationTypeIssue", t, func() {
		convey.Convey("IssueRelationTypeIssue", func() {
			convey.So(IssueRelationTypeIssue, convey.ShouldEqual, 4)
		})
	})
}

//资源
func TestConstIssueRelationTypeResource(t *testing.T) {

	convey.Convey("Test IssueRelationTypeResource", t, func() {
		convey.Convey("IssueRelationTypeResource", func() {
			convey.So(IssueRelationTypeResource, convey.ShouldEqual, 20)
		})
	})
}

//状态
func TestConstIssueRelationTypeStatus(t *testing.T) {

	convey.Convey("Test IssueRelationTypeStatus", t, func() {
		convey.Convey("IssueRelationTypeStatus", func() {
			convey.So(IssueRelationTypeStatus, convey.ShouldEqual, 2)
		})
	})
}

//项目对象类型LangCode
func TestConstProjectObjectTypeLangCodeIteration(t *testing.T) {

	convey.Convey("Test ProjectObjectTypeLangCodeIteration", t, func() {
		convey.Convey("ProjectObjectTypeLangCodeIteration", func() {
			convey.So(ProjectObjectTypeLangCodeIteration, convey.ShouldEqual, "Project.ObjectType.Iteration")
		})
	})
}

func TestConstProjectObjectTypeLangCodeBug(t *testing.T) {

	convey.Convey("Test ProjectObjectTypeLangCodeBug", t, func() {
		convey.Convey("ProjectObjectTypeLangCodeBug", func() {
			convey.So(ProjectObjectTypeLangCodeBug, convey.ShouldEqual, "Project.ObjectType.Bug")
		})
	})
}

func TestConstProjectObjectTypeLangCodeTestTask(t *testing.T) {

	convey.Convey("Test ProjectObjectTypeLangCodeTestTask", t, func() {
		convey.Convey("ProjectObjectTypeLangCodeTestTask", func() {
			convey.So(ProjectObjectTypeLangCodeTestTask, convey.ShouldEqual, "Project.ObjectType.TestTask")
		})
	})
}

func TestConstProjectObjectTypeLangCodeFeature(t *testing.T) {

	convey.Convey("Test ProjectObjectTypeLangCodeFeature", t, func() {
		convey.Convey("ProjectObjectTypeLangCodeFeature", func() {
			convey.So(ProjectObjectTypeLangCodeFeature, convey.ShouldEqual, "Project.ObjectType.Feature")
		})
	})
}

func TestConstProjectObjectTypeLangCodeDemand(t *testing.T) {

	convey.Convey("Test ProjectObjectTypeLangCodeDemand", t, func() {
		convey.Convey("ProjectObjectTypeLangCodeDemand", func() {
			convey.So(ProjectObjectTypeLangCodeDemand, convey.ShouldEqual, "Project.ObjectType.Demand")
		})
	})
}

func TestConstProjectObjectTypeLangCodeTask(t *testing.T) {

	convey.Convey("Test ProjectObjectTypeLangCodeTask", t, func() {
		convey.Convey("ProjectObjectTypeLangCodeTask", func() {
			convey.So(ProjectObjectTypeLangCodeTask, convey.ShouldEqual, "Project.ObjectType.Task")
		})
	})
}

//项目类型LangCode
func TestConstProjectTypeLangCodeNormalTask(t *testing.T) {

	convey.Convey("Test ProjectTypeLangCodeNormalTask", t, func() {
		convey.Convey("ProjectTypeLangCodeNormalTask", func() {
			convey.So(ProjectTypeLangCodeNormalTask, convey.ShouldEqual, "ProjectType.NormalTask")
		})
	})
}

func TestConstProjectTypeLangCodeAgile(t *testing.T) {

	convey.Convey("Test ProjectTypeLangCodeAgile", t, func() {
		convey.Convey("ProjectTypeLangCodeAgile", func() {
			convey.So(ProjectTypeLangCodeAgile, convey.ShouldEqual, "ProjectType.Agile")
		})
	})
}

//流程langCode
func TestConstlangCode(t *testing.T) {

	convey.Convey("Test langCode", t, func() {
		//默认任务流程
		convey.Convey("ProcessLangCodeDefaultTask", func() {
			convey.So(ProcessLangCodeDefaultTask, convey.ShouldEqual, "Process.Issue.DefaultTask")
		})
		//默认敏捷项目任务流程
		convey.Convey("ProcessLangCodeDefaultAgileTask", func() {
			convey.So(ProcessLangCodeDefaultAgileTask, convey.ShouldEqual, "Process.Issue.DefaultAgileTask")
		})
		//默认缺陷流程
		convey.Convey("ProcessLangCodeDefaultBug", func() {
			convey.So(ProcessLangCodeDefaultBug, convey.ShouldEqual, "Process.Issue.DefaultBug")
		})
		//默认测试任务流程
		convey.Convey("ProcessLangCodeDefaultTestTask", func() {
			convey.So(ProcessLangCodeDefaultTestTask, convey.ShouldEqual, "ProcessLangCodeDefaultTestTask")
		})
		//默认项目流程
		convey.Convey("ProcessLangCodeDefaultProject", func() {
			convey.So(ProcessLangCodeDefaultProject, convey.ShouldEqual, "Process.DefaultProject")
		})
		//默认迭代流程
		convey.Convey("ProcessLangCodeDefaultIteration", func() {
			convey.So(ProcessLangCodeDefaultIteration, convey.ShouldEqual, "Process.DefaultIteration")
		})
		//默认特性流程
		convey.Convey("ProcessLangCodeDefaultFeature", func() {
			convey.So(ProcessLangCodeDefaultFeature, convey.ShouldEqual, "Process.Issue.DefaultFeature")
		})
		//默认需求流程
		convey.Convey("ProcessLangCodeDefaultDemand", func() {
			convey.So(ProcessLangCodeDefaultDemand, convey.ShouldEqual, "Process.Issue.DefaultDemand")
		})
	})
}

func TestConstProces(t *testing.T) {

	convey.Convey("Test ObjectType", t, func() {
		//项目流程
		convey.Convey("ProcessPrject", func() {
			convey.So(ProcessPrject, convey.ShouldEqual, 1)
		})
		//迭代流程
		convey.Convey("ProcessIteration", func() {
			convey.So(ProcessIteration, convey.ShouldEqual, 2)
		})
		//问题流程
		convey.Convey("ProcessIssue", func() {
			convey.So(ProcessIssue, convey.ShouldEqual, 3)
		})
	})
}

func TestConstPriorityType(t *testing.T) {

	convey.Convey("Test PriorityType", t, func() {
		//优先级类型-项目
		convey.Convey("PriorityTypeProject", func() {
			convey.So(PriorityTypeProject, convey.ShouldEqual, 1)
		})
		//优先级类型-需求/任务等优先级
		convey.Convey("PriorityTypeIssue", func() {
			convey.So(PriorityTypeIssue, convey.ShouldEqual, 2)
		})
	})
}

func TestConstProcessStatusCategoryProject(t *testing.T) {

	convey.Convey("Test PriorityType", t, func() {
		//项目状态
		convey.Convey("ProcessStatusCategoryProject", func() {
			convey.So(ProcessStatusCategoryProject, convey.ShouldEqual, 1)
		})
		//迭代状态
		convey.Convey("ProcessStatusCategoryIteration", func() {
			convey.So(ProcessStatusCategoryIteration, convey.ShouldEqual, 2)
		})
		//任务状态
		convey.Convey("ProcessStatusCategoryIssue", func() {
			convey.So(ProcessStatusCategoryIssue, convey.ShouldEqual, 3)
		})
	})
}

func TestConstrelationType(t *testing.T) {

	convey.Convey("Test PriorityType", t, func() {
		convey.Convey("UserTeamRelationTypeLeader", func() {
			convey.So(UserTeamRelationTypeLeader, convey.ShouldEqual, 1)
		})
		convey.Convey("UserTeamRelationTypeMember", func() {
			convey.So(UserTeamRelationTypeMember, convey.ShouldEqual, 2)
		})
	})
}

func TestConstReport(t *testing.T) {

	convey.Convey("Test DailyReport", t, func() {
		convey.Convey("DailyReport", func() {
			convey.So(DailyReport, convey.ShouldEqual, 1)
		})
		convey.Convey("WeeklyReport", func() {
			convey.So(WeeklyReport, convey.ShouldEqual, 2)
		})
		convey.Convey("MonthlyReport", func() {
			convey.So(MonthlyReport, convey.ShouldEqual, 3)
		})

	})
}

func TestConstTemplateDirPrefix(t *testing.T) {

	convey.Convey("Test TemplateDirPrefix", t, func() {
		convey.Convey("TemplateDirPrefix", func() {
			convey.So(DailyReport, convey.ShouldEqual, "resources/template/")
		})

	})
}

//资源存储方式
func TestConstResource(t *testing.T) {

	convey.Convey("Test TemplateDirPrefix", t, func() {

		convey.Convey("LocalResource", func() {
			convey.So(LocalResource, convey.ShouldEqual, 1)
		})

		convey.Convey("OssResource", func() {
			convey.So(OssResource, convey.ShouldEqual, 2)
		})

		convey.Convey("DingDiskResource", func() {
			convey.So(DingDiskResource, convey.ShouldEqual, 3)
		})

	})
}

//公用项目
func TestConstProject(t *testing.T) {

	convey.Convey("Test PublicProject", t, func() {

		convey.Convey("PublicProject", func() {
			convey.So(PublicProject, convey.ShouldEqual, 1)
		})

		convey.Convey("PrivateProject", func() {
			convey.So(PrivateProject, convey.ShouldEqual, 2)
		})
	})
}

//公用项目
func TestConstProjectMember(t *testing.T) {

	convey.Convey("Test ProjectMemberEffective", t, func() {

		convey.Convey("ProjectMemberEffective", func() {
			convey.So(ProjectMemberEffective, convey.ShouldEqual, 1)
		})

		convey.Convey("ProjectMemberDisabled", func() {
			convey.So(ProjectMemberDisabled, convey.ShouldEqual, 2)
		})
	})
}
