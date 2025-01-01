import 'package:flutter/material.dart';
import "package:http/src/response.dart";
import 'dart:convert';
import 'dart:developer' as developer;
import 'package:fluttertoast/fluttertoast.dart';

import "../custom_widget.dart";
import "../apiCall/apiRequest.dart";
import '../globalData.dart';

class DeleteWorkflows {
  static void deleteWorkflowToast(String msg) {
    Fluttertoast.showToast(
      msg: msg,
      toastLength: Toast.LENGTH_SHORT,
      timeInSecForIosWeb: 1,
      backgroundColor: Colors.red,
      textColor: Colors.white,
      fontSize: 24.0
    );
  }

  static Future<void> deleteWorkflow(String workflowID) async {
    String apiRoute = "workflows/$workflowID";
    Response response = await ApiRequest.delete(
      apiRoute,
    );
    if (response.statusCode == 504) {
      deleteWorkflowToast(response.body);
      return;
    }
    final responseData = jsonDecode(response.body);
    if (response.statusCode != 200) {
      deleteWorkflowToast('${responseData['error']}');
      developer.log(
          "deleteWorkflow FAIL ERROR : \nError Message : $responseData\nError Code : ${response.statusCode}");
      return;
    }
    deleteWorkflowToast("Workflow deleted");
    Navigator.pop(Globaldata.myContext);
  }

  static Future<void> deleteWorkflowConfirmDialog(String workflowID) {
    return showDialog<void>(
      context: Globaldata.myContext,
      barrierDismissible: false,
      builder: (BuildContext context) {
        return TwoChoiceAlertDialog(
          firstChoiceText: Text(
            "Continue",
            style: TextStyle(
                fontSize: 24, fontWeight: FontWeight.bold, color: Colors.white),
          ),
          secondChoiceText: Text(
            "Cancel",
            style: TextStyle(
              fontSize: 20,
              fontWeight: FontWeight.bold,
              color: Colors.black,
            ),
          ),
          firstChoiceonPressed: () {
            deleteWorkflow(workflowID);
            Navigator.of(context).pop();
          },
          secondChoiceonPressed: () {
            Navigator.of(context).pop();
          },
        );
      },
    );
  }
}