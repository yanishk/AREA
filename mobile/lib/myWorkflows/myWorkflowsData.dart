class MyWorkflowsData {
  static WorkflowsData workflowsData = WorkflowsData.optional();
  static List<WorkflowsData> allWorkflowsData = [];
  static int selectedWorkflowsIndex = -1;
}

class WorkflowsData {
  String id = "";
  String name = "";
  String description = "";
  String serviceId = "";
  String ownerID = "";
  String actionID = "";
  String reactionId = "";
  bool isActivated = false;
  String createdAt = "";
  Map<String, dynamic> actionparam = {};
  Map<String, dynamic> reactionparam = {};

  WorkflowsData.optional([
    this.id = "",
    this.name = "",
    this.ownerID = "",
    this.actionID = "",
    this.reactionId = "",
    this.isActivated = false,
    this.createdAt = "",
    Map<String, dynamic>? actionparam,
    Map<String, dynamic>? reactionparam,
  ]) {
    this.actionparam = actionparam ?? {};
    this.reactionparam = reactionparam ?? {};
  }
}
