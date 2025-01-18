import 'package:flutter/material.dart';
import 'package:auto_size_text/auto_size_text.dart';
import 'package:flutter_switch/flutter_switch.dart';
import 'dart:convert';

import "../custom_widget.dart";
import "../apiCall/apiRequest.dart";
import "myWorkflowsData.dart";
import 'workflowsWidget.dart';
import "myWorkflowsFunctions.dart";
import 'deleteWorkflow.dart';
import "../globalData.dart";

class WorkflowInfosPage extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      theme: ThemeData(
          scaffoldBackgroundColor: const Color.fromARGB(255, 255, 255, 255)),
      home: WorkflowInfosAPICall(),
    );
  }
}

class WorkflowInfosAPICall extends StatefulWidget {
  @override
  WorkflowInfos createState() => WorkflowInfos();
}

class WorkflowInfos extends State<WorkflowInfosAPICall> {
  TextEditingController workflowTitleController = new TextEditingController(text: MyWorkflowsData.allWorkflowsData[MyWorkflowsData.selectedWorkflowsIndex].name);
  String email = Globaldata.userInfos.username;
  Color backgroundColor = HexColor(
    Globaldata.serviceList[MyWorkflowsfunctions.getServiceLocalId(MyWorkflowsData.allWorkflowsData[MyWorkflowsData.selectedWorkflowsIndex].actionID)].color
  );

  bool editNameMode = false;
  String editNameButtonText = "Edit title";
  String modifDate = "";
  String workflowID = MyWorkflowsData.allWorkflowsData[MyWorkflowsData.selectedWorkflowsIndex].id;

  @override
  void initState() {
    super.initState();
    modifDate = MyWorkflowsfunctions.getWorkflowDate();
  }

  Future<void> enableEditName() async {
    setState(() {
      editNameMode = !editNameMode;
      if (editNameMode) {
        editNameButtonText = "Save";
      } else {
        editNameButtonText = "Edit title";
        ApiRequest.put(
          "workflows/$workflowID",
          jsonEncode(<String, dynamic>{
            'name': workflowTitleController.text,
          }),
        );
      }
    });
  }

  Future<void> changeActivationBool() async {
    setState(() {
      MyWorkflowsData.allWorkflowsData[MyWorkflowsData.selectedWorkflowsIndex].isActivated = !MyWorkflowsData.allWorkflowsData[MyWorkflowsData.selectedWorkflowsIndex].isActivated;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SingleChildScrollView(
        child: Column(
          children: [
            WorkflowInfosTitle(
              workflowTitleController: workflowTitleController,
              editTextButtonOnPressed: () {enableEditName();},
              backgroundColor: backgroundColor,
              isEditTextEnable: editNameMode,
              editTextButtonText: editNameButtonText,
              email: email,
            ),
            Padding(padding: EdgeInsets.only(top: 30)),
            Padding(
              padding: EdgeInsets.symmetric(horizontal: 15),
              child: FlutterSwitch(
                showOnOff: true,
                inactiveText: "Connect",
                activeText: "Connected",
                activeTextFontWeight: FontWeight.w900,
                inactiveTextFontWeight: FontWeight.w900,
                inactiveTextColor: Colors.white,
                activeTextColor: Colors.white,
                inactiveColor: const Color.fromARGB(255, 102, 102, 102),
                valueFontSize: 30,
                value: MyWorkflowsData.allWorkflowsData[MyWorkflowsData.selectedWorkflowsIndex].isActivated,
                width: double.infinity,
                height: 70,
                padding: 15,
                activeColor: const Color.fromARGB(255, 34, 34, 34),
                onToggle: (bool) {
                  changeActivationBool();
                  ApiRequest.put(
                    "workflows/$workflowID",
                    jsonEncode(<String, dynamic>{
                      'isactivated': MyWorkflowsData.allWorkflowsData[MyWorkflowsData.selectedWorkflowsIndex].isActivated,
                    }),
                  );
                },
              ),
            ),
            Padding(padding: EdgeInsets.only(top: 30)),
            LeftWidget(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  AutoSizeText(
                    "More details",
                    style: TextStyle(
                      fontSize: 15,
                      fontWeight: FontWeight.bold,
                      color: const Color.fromARGB(255, 102, 102, 102),
                    ),
                    maxLines: 1,
                  ),
                  AutoSizeText(
                    "$modifDate",
                    style: TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.bold,
                      color: Colors.black,
                    ),
                    maxLines: 5,
                  ),
                ],
              )
            ),
            Padding(padding: EdgeInsets.only(top: 30)),
            CustomButton(
              backgroundColor: Color.fromARGB(255, 240, 240, 240),
              text: AutoSizeText(
                "Delete workflow",
                style: TextStyle(
                  fontSize: 30,
                  fontWeight: FontWeight.bold,
                  color: const Color.fromARGB(255, 255, 0, 0)
                ),
                maxLines: 1,
              ),
              onPressed: () {DeleteWorkflows.deleteWorkflowConfirmDialog(workflowID);},
              size: Size(double.infinity, 55),
              padding: EdgeInsets.symmetric(horizontal: 15.0),
            ),
          ],
        ),
      ),
    );
  }
}

class WorkflowInfosTitle extends StatelessWidget {
  const WorkflowInfosTitle({
    super.key,
    required this.workflowTitleController,
    required this.editTextButtonOnPressed,
    this.backgroundColor = Colors.black,
    this.isEditTextEnable = false,
    this.editTextButtonText = "",
    this.email = "",
  });
  final Color backgroundColor;
  final TextEditingController workflowTitleController;
  final bool isEditTextEnable;
  final String editTextButtonText;
  final String email;
  final dynamic editTextButtonOnPressed;

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        color: backgroundColor,
        border: Border.all(
          color: backgroundColor,
          width: 2.0,
        ),
        borderRadius: BorderRadius.circular(4.0),
      ),
      child: InkWell(
        child: Column(
          children: [
            Padding(padding: EdgeInsets.only(top: 30)),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 10),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  GoBackArrow(),
                  EditWheel(),
                ],
              ),
            ),
            Padding(padding: EdgeInsets.only(top: 10)),
            LeftWidget(
                child: Row(
              children: [
                SafeWebImage(
                  url:
                      "${Globaldata.domainName}${
                        Globaldata.serviceList[MyWorkflowsfunctions.getServiceLocalId(MyWorkflowsData.allWorkflowsData[MyWorkflowsData.selectedWorkflowsIndex].actionID)].icon
                      }",
                  width: 50,
                  height: 50,
                ),
                Padding(padding: EdgeInsets.only(left: 10)),
                SafeWebImage(
                  url:
                      "${Globaldata.domainName}${
                        Globaldata.serviceList[MyWorkflowsfunctions.getServiceLocalId(MyWorkflowsData.allWorkflowsData[MyWorkflowsData.selectedWorkflowsIndex].reactionId)].icon
                      }",
                  width: 50,
                  height: 50,
                ),
              ],
            )),
            Padding(padding: EdgeInsets.only(top: 5)),
            WidgetCard(
              backgroundColor: backgroundColor,
              borderRadius: 10,
              child: CustomClassicForm(
                controller: workflowTitleController,
                maxLines: 5,
                hint: '',
                padding:
                    EdgeInsets.symmetric(horizontal: 0, vertical: 0),
                borderWidth: 5,
                fillCollor: backgroundColor,
                enabledBorderColor: Colors.white,
                disabledBorderColor: backgroundColor,
                focusedBorderColor: Colors.white,
                enabled: isEditTextEnable,
                textColor: Colors.white,
              ),
            ),
            LeftWidget(
              child: CustomTextButton(
                text: Text(
                  editTextButtonText,
                  style: TextStyle(
                    fontSize: 20,
                    fontWeight: FontWeight.bold,
                    color: Colors.white,
                  ),
                ),
                onPressed: editTextButtonOnPressed,
                padding: EdgeInsets.zero,
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
    );
  }
}
