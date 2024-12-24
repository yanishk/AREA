import 'package:flutter/material.dart';
import 'package:auto_size_text/auto_size_text.dart';
import "package:http/src/response.dart";
import 'package:flutter_switch/flutter_switch.dart';
import 'dart:convert';
import 'dart:developer' as developer;
import 'package:fluttertoast/fluttertoast.dart';
import "../../apiCall/apiErrorHandling.dart";
import 'package:flutter_phoenix/flutter_phoenix.dart';

import "../custom_widget.dart";
import "../apiCall/apiRequest.dart";
import "myWorkflowsData.dart";
import 'workflowsWidget.dart';
import '../globalData.dart';

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

  static void storeServiceDescription(dynamic responseData, List<Widget> cardList, String workflowsResponseName, Widget workflowsTiles(int index)) {
    if (responseData is Map<String, dynamic> &&
        responseData[workflowsResponseName] is List) {
      List<dynamic> actions = responseData[workflowsResponseName];
      for (int i = 0; !actions.isEmpty; i++) {
        final actionResponse = actions.removeLast();
        if (actionResponse is Map<String, dynamic>) {
          MyWorkflowsData.allWorkflowsData.add(WorkflowsData.optional(
            actionResponse['Id'] ?? '',
            actionResponse['Name'] ?? '',
            actionResponse['OwnerId'] ?? '',
            actionResponse['ActionId'] ?? '',
            actionResponse['ReactionId'] ?? '',
            actionResponse['IsActivated'] ?? false,
            actionResponse['CreatedAt'] ?? '',
            actionResponse['actionparam'] ?? '',
            actionResponse['reactionparam'] ?? '',
          ));
          cardList.add(
            WorkflowCard(
              content: workflowsTiles(i),
            ),
          );
        }
      }
    }
  }
}