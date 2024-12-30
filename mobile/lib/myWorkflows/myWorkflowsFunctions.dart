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

  static Color getServiceColorById (String serviceId) {
    developer.log("Globaldata.actionServicesID: ${Globaldata.actionServicesID}");
    int localActionId = getServiceLocalId(serviceId, Globaldata.actionServicesID);
    int localReactionId = getServiceLocalId(serviceId, Globaldata.reactionServicesID);
    if (localActionId != -1) {
      developer.log("serviceId : $serviceId");
      return HexColor(Globaldata.actionServicesColor[localActionId]);
    }
    if (localReactionId != -1) {
      return HexColor(Globaldata.reactionServicesColor[localReactionId]);
    }
    return Colors.black;
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