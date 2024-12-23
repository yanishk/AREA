import 'package:flutter/material.dart';
import 'package:auto_size_text/auto_size_text.dart';
import "package:http/src/response.dart";
import 'package:flutter_switch/flutter_switch.dart';
import 'dart:convert';

import "../custom_widget.dart";
import "../apiCall/apiRequest.dart";
import "../apiCall/apiErrorHandling.dart";
import "myWorkflowsData.dart";
import 'workflowsWidget.dart';

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
  final String apiRoute = "workflow/create";
  String email = "";
  Color backgroundColor = Color.fromARGB(255, 63, 68, 69);
  bool editNameMode = false;
  String editNameButtonText = "Edit name";
  String modifDate = "";

  @override
  void initState() {
    super.initState();
    getEmail();
    setDateReadable();
  }

  void setDateReadable() {
    List<String> dateSplited = MyWorkflowsData.allWorkflowsData[MyWorkflowsData.selectedWorkflowsIndex].createdAt.split(RegExp(r'[-:T\.]'));
    if (dateSplited.length < 5) {
      return;
    }
    dateSplited.removeAt(dateSplited.length - 1);
    modifDate = "Created on ${dateSplited[1]} ${dateSplited[2]} ${dateSplited[0]} at ${dateSplited[3]}h${dateSplited[4]}";
  }

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

  Future<void> enableEditName() async {
    setState(() {
      editNameMode = !editNameMode;
      if (editNameMode) {
        editNameButtonText = "Save";
      } else {
        editNameButtonText = "Edit name";
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
                  ApiRequest.put(
                    "workflow/modify-activation",
                    jsonEncode(<String, dynamic>{'id': MyWorkflowsData.allWorkflowsData[MyWorkflowsData.selectedWorkflowsIndex].id,}),
                  );
                  changeActivationBool();
                },
              ),
            ),
            Padding(padding: EdgeInsets.only(top: 30)),
            AutoSizeText(
              "$modifDate",
              style: TextStyle(
                fontSize: 20,
                fontWeight: FontWeight.bold,
                color: Colors.black,
              ),
              maxLines: 1,
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
            LeftWidget(
              child: GoBackArrow(),
              padding: const EdgeInsets.symmetric(horizontal: 10),
            ),
            Padding(padding: EdgeInsets.only(top: 10)),
            LeftWidget(
                child: Row(
              children: [
                // SafeWebImage(
                //   url:
                //       "${Globaldata.domainName}${Globaldata.actionServicesIcons[Globaldata.firstService]}",
                //   width: 30,
                //   height: 30,
                // ),
                Padding(padding: EdgeInsets.only(left: 10)),
                // SafeWebImage(
                //   url:
                //       "${Globaldata.domainName}${Globaldata.reactionServicesIcons[Globaldata.secondService]}",
                //   width: 30,
                //   height: 30,
                // ),
              ],
            )),
            Padding(padding: EdgeInsets.only(top: 5)),
            WidgetCard(
              backgroundColor: backgroundColor,
              borderRadius: 10,
              child: CustomClassicForm(
                controller: workflowTitleController,
                maxLines: 5,
                hint: 'Email',
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
                "By $email",
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
