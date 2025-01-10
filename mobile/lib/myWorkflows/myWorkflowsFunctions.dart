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

  static int getServiceLocalId(String uuid) {
    String serviceID = "";
    for (int i = 0; i < Globaldata.actionServicesActions.length && serviceID.isEmpty; i++) {
      if (Globaldata.actionServicesActions[i].id == uuid) {
        serviceID = Globaldata.actionServicesActions[i].serviceid;
      }
    }
    for (int i = 0; i < Globaldata.reactionServicesReaction.length && serviceID.isEmpty; i++) {
      if (Globaldata.reactionServicesReaction[i].id == uuid) {
        serviceID = Globaldata.reactionServicesReaction[i].serviceid;
      }
    }
    if (serviceID.isEmpty) {
      serviceID = uuid;
    }
    for (int i = 0; i < Globaldata.serviceList.length; i++) {
      if (Globaldata.serviceList[i].id == serviceID) {
        return i;
      }
    }
    return 0;
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
              color: HexColor(
                Globaldata.serviceList[getServiceLocalId(MyWorkflowsData.allWorkflowsData[i].actionID)].color
              ),
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