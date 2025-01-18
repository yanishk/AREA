import 'package:flutter/material.dart';
import 'package:auto_size_text/auto_size_text.dart';
import 'dart:convert';
import 'package:fluttertoast/fluttertoast.dart';
import "package:http/src/response.dart";
import 'dart:developer' as developer;

import "../custom_widget.dart";
import "../globalData.dart";
import "../explore/explorePage.dart";
import "../apiCall/apiRequest.dart";
import "../apiCall/apiErrorHandling.dart";
import "createWidget.dart";
import "createWorkflowsData.dart";

class SaveNewWorkflowsPage extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      theme: ThemeData(
          scaffoldBackgroundColor: const Color.fromARGB(255, 255, 255, 255)),
      home: APICall(),
    );
  }
}

class APICall extends StatefulWidget {
  @override
  SaveNewWorkflows createState() => SaveNewWorkflows();
}

class SaveNewWorkflows extends State<APICall> {
  TextEditingController workflowTitleController = new TextEditingController(
      text:
          "If ${CreateData.serviceAction.name}, Then ${CreateData.reactionServiceList[Globaldata.reactionServiceIndex].name} : ${CreateData.serviceReaction.name}");
  final String apiRoute = "workflows";
  String email = Globaldata.userInfos.username;
  String defaultWorkflowName = "";

  @override
  void initState() {
    super.initState();
    developer.log("${CreateData.actionParamListName}");
    developer.log("${CreateData.reactionParamListName}");
    defaultWorkflowName = "If ${CreateData.serviceAction.name} ";
    for (int i = 0; i < CreateData.actionParamListName.length; i++) {
      defaultWorkflowName += "${CreateData.actionParamListName[i]} ";
    }
    defaultWorkflowName += ", Then ${CreateData.reactionServiceList[Globaldata.reactionServiceIndex].name} ${CreateData.serviceReaction.name} ";
    for (int i = 0; i < CreateData.reactionParamListName.length; i++) {
      defaultWorkflowName += "${CreateData.reactionParamListName[i]} ";
    }
    workflowTitleController = new TextEditingController(
      text: defaultWorkflowName,
    );
  }

  Future<void> postCreateWorkflows() async {
    developer.log("${CreateData.serviceAction.parameterList}");
    Response response = await ApiRequest.post(
      apiRoute,
      jsonEncode(<String, dynamic>{
        'actionid': CreateData.serviceAction.id,
        'actionparam': CreateData.serviceAction.parameterList,
        'name': workflowTitleController.value.text.toString(),
        'reactionid': CreateData.serviceReaction.id,
        'reactionparam': CreateData.serviceReaction.parameterList,
      }),
    );
    final responseData = await errorHandler.basicResponseErrorHandler(response, "postCreateWorkflows");
    if (responseData == null) {
      return;
    }
    setState(() {
      Fluttertoast.showToast(
          msg: "Workflows successfully created",
          toastLength: Toast.LENGTH_SHORT,
          backgroundColor: const Color.fromARGB(255, 70, 70, 70),
          textColor: Colors.white,
          fontSize: 24.0);
      Navigator.push(
        Globaldata.myContext,
        MaterialPageRoute(builder: (context) => ExplorePage(),),
      );
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SingleChildScrollView(
        child: Column(
          children: [
            Container(
              decoration: BoxDecoration(
                color: Color.fromARGB(255, 63, 68, 69),
                border: Border.all(
                  color: Color.fromARGB(255, 63, 68, 69),
                  width: 2.0,
                ),
                borderRadius: BorderRadius.circular(4.0),
              ),
              child: InkWell(
                child: Column(
                  children: [
                    Padding(padding: EdgeInsets.only(top: 30)),
                    SelectServiceTopBar(
                        title: "Review and finish", color: Colors.white),
                    Padding(padding: EdgeInsets.only(top: 20)),
                    LeftWidget(
                        child: Row(
                      children: [
                        SafeWebImage(
                          url:
                              "${Globaldata.domainName}${CreateData.serviceAction.icon}",
                          width: 30,
                          height: 30,
                        ),
                        Padding(padding: EdgeInsets.only(left: 10)),
                        SafeWebImage(
                          url:
                              "${Globaldata.domainName}${CreateData.serviceReaction.icon}",
                          width: 30,
                          height: 30,
                        ),
                      ],
                    )),
                    Padding(padding: EdgeInsets.only(top: 5)),
                    LeftWidget(
                      child: AutoSizeText(
                        "Workflow title",
                        style: TextStyle(
                          fontSize: 20,
                          fontWeight: FontWeight.bold,
                          color: Colors.white,
                        ),
                        maxLines: 1,
                      ),
                    ),
                    WidgetCard(
                      backgroundColor: Colors.white,
                      borderRadius: 10,
                      child: CustomClassicForm(
                        controller: workflowTitleController,
                        maxLines: 5,
                        hint: 'Workflow title',
                        padding:
                            EdgeInsets.symmetric(horizontal: 0, vertical: 0),
                        focusedBorderColor: Colors.white,
                        borderWidth: 0,
                      ),
                    ),
                    LeftWidget(
                      child: AutoSizeText(
                        "by $email",
                        style: TextStyle(
                          fontSize: 20,
                          fontWeight: FontWeight.bold,
                          color: Colors.white,
                        ),
                        maxLines: 1,
                      ),
                    ),
                    Padding(padding: EdgeInsets.only(top: 10)),
                  ],
                ),
              ),
            ),
            Padding(padding: EdgeInsets.only(top: 30)),
            CustomButton(
              backgroundColor: Color.fromARGB(255, 34, 34, 34),
              text: AutoSizeText(
                "Finish",
                style: TextStyle(
                    fontSize: 30,
                    fontWeight: FontWeight.bold,
                    color: Colors.white),
              ),
              onPressed: () {
                postCreateWorkflows();
              },
              size: Size(double.infinity, 60),
              padding: EdgeInsets.symmetric(horizontal: 15.0),
            ),
          ],
        ),
      ),
    );
  }
}
