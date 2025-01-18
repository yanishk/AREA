import 'package:flutter/material.dart';
import 'package:auto_size_text/auto_size_text.dart';
import 'dart:developer' as developer;
import 'package:fluttertoast/fluttertoast.dart';

import "../../custom_widget.dart";
import "../../globalData.dart";
import '../createWidget.dart';
import '../createWorkflowsData.dart';
import "../../auth/webviewAuth/connectServiceAuth.dart";
import 'actionOrReactionsFunctions.dart';

bool hasToLogin = true;

class SelectParamActionOrReactionPage extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    hasToLogin = CreateData.hasToLogin;
    return MaterialApp(
      theme: ThemeData(
          scaffoldBackgroundColor: HexColor(CreateData.actualServiceList[CreateData.selectedServiceId].color)),
      home: SelectParamActionOrReactionAPICall(),
    );
  }
}

class SelectParamActionOrReactionAPICall extends StatefulWidget {
  @override
  SelectParamActionOrReaction createState() => SelectParamActionOrReaction();
}

class SelectParamActionOrReaction extends State<SelectParamActionOrReactionAPICall> {
  String apiRoute = "actions";
  List<Widget> parameterFiledList = [];
  Map<int, String> selectedItems = {};
  Map<int, TextEditingController> textInputController = {};
  Map<int, String> selectedParameterId = {};
  String asanaWorkspaceId = "";
  List<int> textInputControllerIndexList = [];
  List<String> actionOrReactionParamListName = [];
  Map<int, String> actionOrReactionParamStringMap = {};

  @override
  void initState() {
    super.initState();
    if (!isAnAction) {
      setState(() {
        apiRoute = "reactions";
      });
    }
    if (CreateData.selectedActionOrReactionData.parameter.isEmpty && !hasToLogin) {
      if (isAnAction) {
        setState(() {
          CreateData.actionParamListName = [];
        });
      } else {
        setState(() {
          CreateData.reactionParamListName = [];
        });
      }
      CreateData.selectOrUpdateWorkflow = true;
      Navigator.pop(Globaldata.myContext);
    }
    fillParameterFiledList();
  }

  void fillParameterFiledList() {
    parameterFiledList = [];
    for (int i = 0; i < CreateData.selectedActionOrReactionData.parameter.length; i++) {
      parameterFiledList.add(parameterFieldWidget(i));
    }
  }

  void createDropdownFromRoute(int index, dynamic response, String listName) {
    List<ActionOrReactionsParameter> parameterList = CreateData.selectedActionOrReactionData.parameter;
    List<dynamic> responseList = response[listName];
    List<String> nameList = [];
    List<String> idList = [];
    for (int i = 0; i < responseList.length; i++) {
      String itemName = responseList[i].toString().split("name: ").last.split(",").first;
      if (!nameList.contains(itemName.split("}").first)) {
        nameList.add(itemName.split("}").first);
      }
      if (parameterList[index].name == "channel") {
        idList.add(responseList[i]["id"]);
      } else {
        String itemId = responseList[i].toString().split("id: ").last.split(",").first;
        idList.add(itemId.split("}").first);
      }
    }
    if (nameList.length < 1) {
      developer.log("No : ${parameterList[index].name}");
      return;
    }
    if (!selectedItems.containsKey(index)) {
      selectedItems[index] = nameList.first;
      actionOrReactionParamStringMap[index] = "${selectedItems[index]}";
      selectedParameterId[index] = idList.first;
      int selectedParamIndex = nameList.indexOf(selectedItems[index] ?? "");
      if (selectedParameterId[index]?[0] != "{") {
        CreateData.selectedActionOrReactionData.parameterList[parameterList[index].name] = selectedParameterId[index];
      } else {
        CreateData.selectedActionOrReactionData.parameterList[parameterList[index].name] = nameList[selectedParamIndex];
      }
      if (parameterList[index].name == "workspace") {
        asanaWorkspaceId = idList.first;
      }
    }
    setState(() {
      parameterFiledList.add(
        CreateDropdown(
          index: index,
          onChanged: (String? selectedValue) {
            setState(() {
              selectedItems[index] = selectedValue ?? "";
              selectedParameterId[index] = idList[nameList.indexOf(selectedValue ?? "")];
              int selectedParamIndex = nameList.indexOf(selectedItems[index] ?? "");
              actionOrReactionParamStringMap[index] = "$selectedValue";
              if (selectedParameterId[index]?[0] != "{") {
                developer.log("selectedParameterId[index] : ${selectedParameterId[index]}");
                CreateData.selectedActionOrReactionData.parameterList[parameterList[index].name] = selectedParameterId[index];
              } else {
                developer.log("nameList[selectedParamIndex] : ${nameList[selectedParamIndex]}");
                CreateData.selectedActionOrReactionData.parameterList[parameterList[index].name] = nameList[selectedParamIndex];
              }
              if (parameterList[index].name == "workspace") {
                asanaWorkspaceId = idList[nameList.indexOf(selectedValue ?? "")];
                selectedItems.clear();
                selectedItems[index] = selectedValue ?? "";
              }
              fillParameterFiledList();
            });
          },
          stringList: nameList,
          selectedItems: selectedItems,
          dropdownName: "${parameterList[index].name} : "
        )
      );
    });
  }

  Future<void> getRouteData(String route, int index) async {
    if (CreateData.actualServiceList[CreateData.selectedServiceId].name == "Asana") {
      route = "$route?id=$asanaWorkspaceId";
    } else if (selectedParameterId[index - 1] != null) {
      route = "$route?id=${selectedParameterId[index - 1]}";
    }
    final response = await ActionsOrReactionsFunctions.getRouteResponse(route);
    if (response == null) {
      return;
    }
    String listName = response.toString().split(":").first.split("{").last;
    if (!(response[listName] is List<dynamic>)) {
      developer.log("new response : ${ActionsOrReactionsFunctions.getMapKeyName(response)}");
      dynamic mapKeyName = ActionsOrReactionsFunctions.getMapKeyName(response);
      if (mapKeyName == null) {
        return;
      }
      listName = mapKeyName.toString().split(":").first.split("{").last;
      createDropdownFromRoute(index, mapKeyName, listName);
      return;
    }
    createDropdownFromRoute(index, response, listName);
  }

  List<String> fillDropDownStringList(List<ActionOrReactionsParameter> parameterList, int index) {
    if (parameterList[index].values.isEmpty) {
      if (parameterList[index].route.isEmpty) {
        return [];
      }
      getRouteData(parameterList[index].route, index);
      return [];
    }
    if (parameterList[index].isexhaustive) {
      return parameterList[index].values;
    }
    List<String> splitValuesRange = parameterList[index].values.first.split("-");
    int valueRangeStart = int.parse(splitValuesRange.first);
    int valueRangeEnd = int.parse(splitValuesRange.last);
    List<String> finalStringList = [];
    for (int i = valueRangeStart; i <= valueRangeEnd; i++) {
      finalStringList.add(i.toString());
    }
    return finalStringList;
  }

  Widget parameterFieldWidget(int index) {
    List<ActionOrReactionsParameter> parameterList = CreateData.selectedActionOrReactionData.parameter;
    List<String> dropdownStringList = fillDropDownStringList(parameterList, index);
    if (dropdownStringList.isNotEmpty) {
      if (!selectedItems.containsKey(index)) {
        selectedItems[index] = dropdownStringList.first;
        int selectedParamIndex = dropdownStringList.indexOf(selectedItems[index] ?? "") + 1;
        actionOrReactionParamStringMap[index] = "${selectedItems[index]}";
        if (dropdownStringList.first == "0") {
          selectedParamIndex--;
        }
        if (parameterList[index].type == "string") {
          CreateData.selectedActionOrReactionData.parameterList[parameterList[index].name] = dropdownStringList.first;
        } else {
          CreateData.selectedActionOrReactionData.parameterList[parameterList[index].name] = selectedParamIndex;
        }
      }
      return CreateDropdown(
        index: index,
        onChanged: (String? selectedValue) {
          setState(() {
            selectedItems[index] = selectedValue ?? "";
            int selectedParamIndex = dropdownStringList.indexOf(selectedItems[index] ?? "") + 1;
            actionOrReactionParamStringMap[index] = "$selectedValue";
            if (dropdownStringList.first == "0") {
              selectedParamIndex--;
            }
            if (parameterList[index].type == "string") {
              CreateData.selectedActionOrReactionData.parameterList[parameterList[index].name] = selectedValue;
            } else {
              CreateData.selectedActionOrReactionData.parameterList[parameterList[index].name] = selectedParamIndex;
            }
            fillParameterFiledList();
          });
        },
        stringList: dropdownStringList,
        selectedItems: selectedItems,
        dropdownName: "${parameterList[index].name} : "
      );
    }
    if (parameterList[index].route.isEmpty) {
      if (textInputController[index] == null) {
        textInputController[index] = new TextEditingController(text: "");
      }
      textInputControllerIndexList.add(index);
      return Padding(
        padding: EdgeInsets.only(top: 20),
        child: Container(
          padding: EdgeInsets.all(0.0),
          constraints: BoxConstraints(maxHeight: 200.0, maxWidth: (MediaQuery.of(Globaldata.myContext).size.width - 40)),
          child: SingleChildScrollView(
            child: CustomClassicForm(
              hint: "${parameterList[index].name} : ",
              controller: textInputController[index],
              maxLines: 5,
              padding:
                  EdgeInsets.symmetric(horizontal: 0, vertical: 0),
              focusedBorderColor: Colors.white,
              borderWidth: 0,
            ),
          ),
        ),
      );
    }
    return Column();
  }

  Future<void> handleServiceAuthButton() async {
    final authResponse = await ConnectServiceAuthWebview.connectService("service", CreateData.actualServiceList[CreateData.selectedServiceId].name);
    if (!authResponse) {
      return;
    }
    setState(() {
      hasToLogin = false;
    });
    if (CreateData.selectedActionOrReactionData.parameter.isEmpty) {
      CreateData.selectOrUpdateWorkflow = true;
      Navigator.pop(Globaldata.myContext);
    }
    fillParameterFiledList();
  }

  Widget connectButton() {
    Color buttonColor = HexColor(CreateData.actualServiceList[CreateData.selectedServiceId].color);
    int r = buttonColor.red ~/ 2;
    int g = buttonColor.green ~/ 2;
    int b = buttonColor.blue ~/ 2;
    return CustomButton(
      backgroundColor: Color.fromARGB(255, r, g, b),
      text: AutoSizeText(
        "Connect",
        style: TextStyle(
            fontSize: 30,
            fontWeight: FontWeight.bold,
            color: Colors.white),
      ),
      onPressed: () {
        handleServiceAuthButton();
      },
      size: Size(double.infinity, 60),
      padding: EdgeInsets.symmetric(horizontal: 15.0),
    );
  }

  @override
  Widget build(BuildContext context) {
    if (hasToLogin) {
      return Scaffold(
        body: SingleChildScrollView(
          child: Column(
            children: [
              Padding(padding: EdgeInsets.only(top: 30)),
              SelectActionOrReactionTopBar(title: "Connect service", color: Colors.white),
              Padding(padding: EdgeInsets.only(top: 75)),
              SafeWebImage(
                url: "${Globaldata.domainName}${CreateData.actualServiceList[CreateData.selectedServiceId].icon}",
                width: 120,
                height: 120,
              ),
              AutoSizeText(
                "Connect to ${CreateData.actualServiceList[CreateData.selectedServiceId].name} to continue",
                style: TextStyle(
                  fontSize: 20,
                  fontWeight: FontWeight.bold,
                  color: Colors.white
                ),
              ),
              Padding(padding: EdgeInsets.only(top: 50)),
              connectButton(),
            ],
          ),
        ),
      );
    }
    return Scaffold(
      body: SingleChildScrollView(
        child: Column(
          children: [
            Padding(padding: EdgeInsets.only(top: 30)),
            SelectActionOrReactionTopBar(title: "Complete $apiRoute", color: Colors.white),
            Padding(padding: EdgeInsets.only(top: 50)),
            CompleteActionOrReactionTop(),
            Padding(padding: EdgeInsets.only(top: 15)),
            LeftWidget(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: parameterFiledList,
              ),
            ),
            Padding(padding: EdgeInsets.only(top: 15)),
            CustomButton(
              backgroundColor: ActionsOrReactionsFunctions.fadeColor(HexColor(CreateData.actualServiceList[CreateData.selectedServiceId].color), 2),
              text: AutoSizeText(
                "Continue",
                style: TextStyle(
                    fontSize: 25,
                    fontWeight: FontWeight.bold,
                    color: Colors.white),
              ),
              onPressed: () {
                if (isAnAction) {
                  setState(() {
                    CreateData.actionParamListName = [];
                    actionOrReactionParamListName = CreateData.actionParamListName;
                  });
                } else {
                  setState(() {
                    apiRoute = "reactions";
                    CreateData.reactionParamListName = [];
                    actionOrReactionParamListName = CreateData.reactionParamListName;
                  });
                }
                Map<String, dynamic> parameterList = CreateData.selectedActionOrReactionData.parameterList;
                for (int i = 0; i < textInputControllerIndexList.length; i++) {
                  int controllerIndex = textInputControllerIndexList[i];
                  if (textInputController[controllerIndex] != null) {
                    String controlerText = textInputController[controllerIndex]!.text;
                    ActionOrReactionsParameter thisParameter = CreateData.selectedActionOrReactionData.parameter[controllerIndex];
                    if (thisParameter.type == "int") {
                      if (controlerText.isEmpty) {
                        Fluttertoast.showToast(
                          msg: "Parameter ${thisParameter.name} is missing",
                          toastLength: Toast.LENGTH_SHORT,
                          backgroundColor: const Color.fromARGB(255, 70, 70, 70),
                          textColor: Colors.white,
                          fontSize: 24.0
                        );
                        return;
                      }
                      try {
                        parameterList[thisParameter.name] = int.parse(controlerText);
                      } catch (e) {
                        Fluttertoast.showToast(
                          msg: "Parameter ${thisParameter.name} is not a number",
                          toastLength: Toast.LENGTH_SHORT,
                          backgroundColor: const Color.fromARGB(255, 70, 70, 70),
                          textColor: Colors.white,
                          fontSize: 24.0
                        );
                        return;
                      }
                    } else {
                      parameterList[thisParameter.name] = controlerText;
                    }
                    actionOrReactionParamStringMap[controllerIndex] = "${thisParameter.name} : $controlerText";
                  }
                }
                developer.log("parameter : ${parameterList}");
                for (int i = 0; i < actionOrReactionParamStringMap.length; i++) {
                  actionOrReactionParamListName.add(actionOrReactionParamStringMap[i] ?? "");
                }
                CreateData.selectOrUpdateWorkflow = true;
                Navigator.pop(Globaldata.myContext);
              },
              size: Size(double.infinity, 70),
              padding: EdgeInsets.symmetric(horizontal: 15.0),
            ),
          ],
        ),
      ),
    );
  }
}

class CompleteActionOrReactionTop extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        SafeWebImage(
          url:
              "${Globaldata.domainName}${CreateData.actualServiceList[CreateData.selectedServiceId].icon}",
          width: 75,
          height: 75,
        ),
        Padding(padding: EdgeInsets.only(top: 10)),
        AutoSizeText(
          CreateData.selectedActionOrReactionData.name,
          style: TextStyle(
            fontSize: 25,
            fontWeight: FontWeight.bold,
            color: Colors.white,
          ),
          maxLines: 2,
        ),
        Padding(padding: EdgeInsets.only(top: 5)),
        Padding(
          padding: EdgeInsets.symmetric(horizontal: 10.0),
          child: AutoSizeText(
            CreateData.selectedActionOrReactionData.description,
            style: TextStyle(
              fontSize: 15,
              fontWeight: FontWeight.bold,
              color: Colors.white,
            ),
            textAlign: TextAlign.center,
            maxLines: 10,
          ),
        ),
      ],
    );
  }
}

class CreateDropdown extends StatelessWidget {
  const CreateDropdown({
    super.key,
    required this.onChanged,
    required this.stringList,
    required this.selectedItems,
    required this.index,
    this.dropdownName = "",
  });

  final void Function(String?)? onChanged;
  final List<String> stringList;
  final Map<int, String> selectedItems;
  final int index;
  final String dropdownName;

  @override
  Widget build(BuildContext context) {
    List<String> finalStringList = [];
    for (int i = 0; i < stringList.length; i++) {
      if (!finalStringList.contains(stringList[i])) {
        finalStringList.add(stringList[i]);
      }
    }
    if (selectedItems[index] == null) {
      return Column();
    }
    return Container(
      padding: EdgeInsets.all(0.0),
      constraints: BoxConstraints(maxHeight: 100.0, maxWidth: (MediaQuery.of(Globaldata.myContext).size.width - 40)),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Padding(padding: EdgeInsets.only(top: 20)),
          AutoSizeText(
            dropdownName,
            style: TextStyle(
              fontSize: 20,
              fontWeight: FontWeight.bold,
              color: Colors.white,
            ),
            textAlign: TextAlign.center,
            maxLines: 10,
          ),
          DropdownButton<String>(
            isExpanded: true,
            value: selectedItems[index] ?? "No data",
            hint: AutoSizeText(
              selectedItems[index] ?? "No data",
              style: TextStyle(
                fontSize: 20,
                fontWeight: FontWeight.bold,
                color: Colors.white,
              ),
              maxLines: 1,
            ),
            items: finalStringList.map((String item) {
              return DropdownMenuItem<String>(
                value: item,
                child: AutoSizeText(
                  item,
                  maxLines: 1,
                ),
              );
            }).toList(),
            onChanged: onChanged,
            dropdownColor: Colors.grey[800],
            style: TextStyle(
              color: Colors.white,
              fontSize: 20,
              fontWeight: FontWeight.bold,
            ),
          )
        ]
      ),
    );
  }
}
