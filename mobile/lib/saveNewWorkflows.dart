import 'package:flutter/material.dart';
import 'package:auto_size_text/auto_size_text.dart';
import 'dart:developer' as developer;
import 'dart:convert';
import 'package:fluttertoast/fluttertoast.dart';
import "package:http/src/response.dart";

import "custom_widget.dart";
import "globalData.dart";
import "explore/explorePage.dart";
import "apiCall/apiRequest.dart";

class SaveNewWorkflowsPage extends StatelessWidget {
  const SaveNewWorkflowsPage({
    super.key,
    this.isAnAction = true,
  });
  final bool isAnAction;

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      theme: ThemeData(
          scaffoldBackgroundColor: const Color.fromARGB(255, 255, 255, 255)),
      home: APICall(isAnAction: isAnAction),
    );
  }
}

class APICall extends StatefulWidget {
  const APICall({
    super.key,
    this.isAnAction = true,
  });
  final bool isAnAction;

  @override
  SaveNewWorkflows createState() => SaveNewWorkflows(isAnAction: isAnAction);
}

class SaveNewWorkflows extends State<APICall> {
  SaveNewWorkflows({
    required this.isAnAction,
  });
  final bool isAnAction;
  TextEditingController workflowTitleController = new TextEditingController(
      text:
          "If ${Globaldata.serviceAction.name} ${Globaldata.serviceAction.parameter}, Then ${Globaldata.reactionServicesName[Globaldata.secondService]} : ${Globaldata.serviceReaction.name}");
  final String apiRoute = "workflow/create";
  String email = "";
  bool isInit = true;

  Future<void> postCreateWorkflows() async {
    Response response = await ApiRequest.post(
      apiRoute,
      jsonEncode(<String, dynamic>{
        'actionid': Globaldata.serviceAction.id,
        'actionparam': Globaldata.serviceAction.parameterList,
        'name': workflowTitleController.value.text.toString(),
        'reactionid': Globaldata.serviceReaction.id,
        'reactionparam': Globaldata.serviceReaction.parameterList,
      }),
    );
    if (response.statusCode == 504) {
      developer.log('${response.body}');
      return;
    }
    final responseData = jsonDecode(response.body);
    if (response.statusCode != 200) {
      developer.log(
          "postCreateWorkflows FAIL ERROR : \nError Message : $responseData\nError Code : ${response.statusCode}");
      return;
    }
    setState(() {
      Fluttertoast.showToast(
          msg: "Workflows successfully created",
          toastLength: Toast.LENGTH_SHORT,
          timeInSecForIosWeb: 1,
          backgroundColor: const Color.fromARGB(255, 70, 70, 70),
          textColor: Colors.white,
          fontSize: 24.0);
      Navigator.push(
        Globaldata.myContext,
        MaterialPageRoute(builder: (context) => ExplorePage(),),
      );
    });
  }

  Future<void> getEmail() async {
    Response response = await ApiRequest.get("user/email");
    if (response.statusCode == 504) {
      developer.log('${response.body}');
      return;
    }
    final responseData = jsonDecode(response.body);
    if (response.statusCode != 200) {
      developer.log(
          "getEmail FAIL ERROR : \nError Message : $responseData\nError Code : ${response.statusCode}");
      return;
    }
    setState(() {
      email = responseData["email"];
    });
  }

  @override
  Widget build(BuildContext context) {
    if (isInit == true) {
      isInit = false;
      getEmail();
    }
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
                              "${Globaldata.domainName}${Globaldata.actionServicesIcons[Globaldata.firstService]}",
                          width: 30,
                          height: 30,
                        ),
                        Padding(padding: EdgeInsets.only(left: 10)),
                        SafeWebImage(
                          url:
                              "${Globaldata.domainName}${Globaldata.reactionServicesIcons[Globaldata.secondService]}",
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
                        hint: 'Email',
                        padding:
                            EdgeInsets.symmetric(horizontal: 0, vertical: 0),
                        focusedBorderColor: Colors.white,
                        borderWidth: 0,
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

class SelectServiceTopBar extends StatelessWidget {
  const SelectServiceTopBar({
    super.key,
    required this.title,
    this.color = Colors.black,
  });

  final String title;
  final Color color;

  @override
  Widget build(BuildContext context) {
    return Stack(
      alignment: Alignment.center,
      children: [
        Positioned(
          left: 0,
          child: IconButton(
            onPressed: () {
              Navigator.pop(Globaldata.myContext);
            },
            icon: Icon(
              Icons.arrow_back,
              size: 32,
              color: color,
            ),
          ),
        ),
        Center(
          child: AutoSizeText(
            title,
            style: TextStyle(
              fontSize: 30,
              fontWeight: FontWeight.bold,
              color: color,
            ),
            maxLines: 1,
          ),
        ),
      ],
    );
  }
}
