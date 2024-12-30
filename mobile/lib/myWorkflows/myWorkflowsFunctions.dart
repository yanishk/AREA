import 'package:flutter/material.dart';
import "package:http/src/response.dart";
import 'package:namer_app/globalData.dart';
import 'dart:convert';
import 'dart:developer' as developer;

import "../apiCall/apiRequest.dart";
import "myWorkflowsData.dart";
import 'workflowsWidget.dart';
import "../custom_widget.dart";

class MyWorkflowsfunctions {
  static String getWorkflowDate() {
    List<String> dateSplited = MyWorkflowsData.allWorkflowsData[MyWorkflowsData.selectedWorkflowsIndex].createdAt.split(RegExp(r'[-:T\.]'));
    if (dateSplited.length < 5) {
      return "";
    }
    dateSplited.removeAt(dateSplited.length - 1);
    return "Created on ${dateSplited[1]} ${dateSplited[2]} ${dateSplited[0]}\nAt ${dateSplited[3]}h${dateSplited[4]}";
  }

  static Future<dynamic> getWorkflowsResponse(String route) async {
    Response response = await ApiRequest.get(route);
    if (response.statusCode == 504) {
      developer.log('${response.body}');
      developer.log('${response.body}');
      return null;
    }
    final responseData = jsonDecode(response.body);
    if (response.statusCode != 200) {
      developer.log(
          "getWorkflowsResponse FAIL ERROR : \nError Message : $responseData\nError Code : ${response.statusCode}");
      return null;
    }
    return responseData;
  }

  static int getServiceLocalId(String serviceId, List<String> seviceIdList) {
    for (int i = 0; i < seviceIdList.length; i++) {
      if (seviceIdList[i] == serviceId) {
        return i;
      }
    }
    return -1;
  }

  static String getServiceIdWithActionId(String actionId) {
    String serviceID = "";
    for (int i = 0; i < Globaldata.actionServicesActions.length; i++) {
      if (Globaldata.actionServicesActions[i].id == actionId) {
        serviceID = Globaldata.actionServicesActions[i].serviceid;
      }
    }
    if (serviceID != "") {
      return serviceID;
    }
    for (int i = 0; i < Globaldata.reactionServicesReaction.length; i++) {
      if (Globaldata.reactionServicesReaction[i].id == actionId) {
        serviceID = Globaldata.reactionServicesReaction[i].serviceid;
      }
    }
    return serviceID;
  }

  static Color getServiceColorById (String serviceId) {
    int localActionId = getServiceLocalId(getServiceIdWithActionId(serviceId), Globaldata.actionServicesID);
    if (localActionId != -1) {
      return HexColor(Globaldata.actionServicesColor[localActionId]);
    }
    return Colors.black;
  }

  static String getServiceLogoRoute (String serviceId) {
    int localActionId = getServiceLocalId(getServiceIdWithActionId(serviceId), Globaldata.actionServicesID);
    int localReactionId = getServiceLocalId(getServiceIdWithActionId(serviceId), Globaldata.reactionServicesID);
    if (localActionId != -1) {
      return Globaldata.actionServicesIcons[localActionId];
    }
    if (localReactionId != -1) {
      return Globaldata.reactionServicesIcons[localReactionId];
    }
    return "";
  }

  static void storeServiceDescription(dynamic responseData, List<Widget> cardList, String workflowsResponseName, Future<void> onTapFunction(int index), String email) {
    if (responseData is Map<String, dynamic> &&
        responseData[workflowsResponseName] is List) {
      List<dynamic> actions = responseData[workflowsResponseName];
      for (int i = 0; !actions.isEmpty; i++) {
        final actionResponse = actions.removeLast();
        if (actionResponse is Map<String, dynamic>) {
          MyWorkflowsData.allWorkflowsData.add(WorkflowsData.optional(
            actionResponse['id'] ?? '',
            actionResponse['name'] ?? '',
            actionResponse['ownerid'] ?? '',
            actionResponse['actionid'] ?? '',
            actionResponse['reactionid'] ?? '',
            actionResponse['isactivated'] ?? false,
            actionResponse['createdat'] ?? '',
            actionResponse['actionparam'] ?? '',
            actionResponse['reactionparam'] ?? '',
          ));
          cardList.add(
            WorkflowCard(
              color: getServiceColorById(MyWorkflowsData.allWorkflowsData[i].actionID),
              content: WorkflowsTiles(
                index: i,
                onTap: () {onTapFunction(i);},
                email: email,
              )
            ),
          );
        }
      }
    }
  }
}