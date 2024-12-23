import 'package:flutter/material.dart';
import 'package:auto_size_text/auto_size_text.dart';
import "package:http/src/response.dart";
import 'dart:convert';
import 'dart:developer' as developer;
import 'package:flutter_switch/flutter_switch.dart';

import "../apiCall/apiRequest.dart";
import "../apiCall/apiErrorHandling.dart";

import "../custom_widget.dart";
import "../footer_page.dart";
import "myWorkflowsData.dart";
import "workflowsWidget.dart";
import "../globalData.dart";
import 'workflowInfos.dart';

class DisplayWorkflows extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      theme: new ThemeData(
          scaffoldBackgroundColor: const Color.fromARGB(255, 255, 255, 255)),
      home: DisplayWorkflowsStateful(),
    );
  }
}

class DisplayWorkflowsStateful extends StatefulWidget {
  @override
  DisplayWorkflowsMain createState() => DisplayWorkflowsMain();
}

class DisplayWorkflowsMain extends State<DisplayWorkflowsStateful> {
  @override
  void initState() {
    super.initState();
    getEmail();
    getWorkflows();
  }

  String apiRoute = "workflow/retrieve";
  String workflowsResponseName = "workflows";
  List<Widget> workflowsCardList = [];
  String email = "";

  Future<void> getEmail() async {
    Response response = await ApiRequest.get("user/email");
    final responseData = await errorHandler.basicResponseErrorHandler(response, "getEmail");
    if (responseData == null) {
      return;
    }
    setState(() {
      email = responseData["email"];
    });
  }

  Widget workflowsTiles(int index) {
    return InkWell(
      onTap: () {
        MyWorkflowsData.selectedWorkflowsIndex = index;
        Navigator.push(
          Globaldata.myContext,
          MaterialPageRoute(builder: (context) => WorkflowInfosPage(),),
        );
      },
      child: Column (
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            MyWorkflowsData.allWorkflowsData[index].name,
            style: TextStyle(
              fontSize: 25,
              fontWeight: FontWeight.bold,
              color: Colors.white,
            ),
          ),
          DoubleWidget(
            padding: const EdgeInsets.only(left: 0),
            mainAxisAlignment: MainAxisAlignment.start,
            firstChild: Text(
              "By ",
              style: TextStyle(
                fontSize: 15,
                fontWeight: FontWeight.bold,
                color: const Color.fromARGB(255, 179, 179, 179),
              ),
            ),
            secondChild: Text(
              "$email",
              style: TextStyle(
                fontSize: 15,
                fontWeight: FontWeight.bold,
                color: Colors.white,
              ),
            ),
          ),

          Padding(padding: EdgeInsets.only(top: 5),),
          LeftWidget(
            padding: const EdgeInsets.symmetric(horizontal: 0),
            child: FlutterSwitch(
              showOnOff: true,
              inactiveText: "Connect",
              activeText: "Connected",
              activeTextFontWeight: FontWeight.normal,
              inactiveTextFontWeight: FontWeight.normal,
              inactiveTextColor: Colors.white,
              activeTextColor: Colors.white,
              inactiveColor: const Color.fromARGB(255, 102, 102, 102),
              value: MyWorkflowsData.allWorkflowsData[index].isActivated,
              width: MediaQuery.sizeOf(Globaldata.myContext).width / 3,
              activeColor: const Color.fromARGB(255, 0, 132, 255),
              onToggle: (bool) {},
            ),
          ),
          Padding(padding: EdgeInsets.only(top: 10)),
        ],
      ),
    );
  }

  void storeServiceDescription(dynamic responseData) {
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
          workflowsCardList.add(
            WorkflowCard(
              content: workflowsTiles(i),
            ),
          );
        }
      }
    }
  }

  Future<dynamic> getWorkflowsResponse(String route) async {
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

  Future<void> getWorkflows() async {
    final response = await getWorkflowsResponse(apiRoute);
    if (response == null) {
      return;
    }
    setState(() {
      storeServiceDescription(response);
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SingleChildScrollView(
        child: Column(
          children: [
            Padding(padding: EdgeInsets.only(top: 30)),
            AutoSizeText("My workflows",
              style: TextStyle(fontSize: 30,fontWeight: FontWeight.bold,color: Colors.black,),
              maxLines: 1,
            ),
            Column(
              children: workflowsCardList,
            ),
          ],
        ),
      ),
      bottomNavigationBar: FooterNavigationBar(currentIndex: 0),
    );
  }
}
