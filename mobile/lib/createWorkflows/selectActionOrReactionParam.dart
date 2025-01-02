import 'package:flutter/material.dart';
import 'package:auto_size_text/auto_size_text.dart';
import 'package:flutter/services.dart';
import 'dart:developer' as developer;

import "../custom_widget.dart";
import "../globalData.dart";
import 'createWidget.dart';
import 'createWorkflowsData.dart';

class SelectParamActionOrReactionPage extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      theme: ThemeData(
          scaffoldBackgroundColor: HexColor(CreateData.serviceColorList[CreateData.selectedServiceId])),
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
  Map<int, TextEditingController> textInputControler = {};

  @override
  void initState() {
    super.initState();
    if (!isAnAction) {
      setState(() {
        apiRoute = "reactions";
      });
    }
    fillParameterFiledList();
  }

  void fillParameterFiledList() {
    parameterFiledList = [];
    for (int i = 0; i < CreateData.selectedActionOrReactionData.parameter.length; i++) {
      parameterFiledList.add(parameterFieldWidget(i));
    }
  }

  List<String> fillDropDownStringList(List<ActionOrReactionsParameter> parameterList, int index) {
    if (parameterList[index].values.isEmpty) {
      return [];
    }
    if (parameterList[index].isexhaustive) {
      return parameterList[index].values;
    }
    List<String> splitValuesRange = parameterList[index].values.first.split("-");
    int firstValue = int.parse(splitValuesRange.first);
    int secondValue = int.parse(splitValuesRange.last);
    List<String> finalStringList = [];
    for (int i = firstValue; i <= secondValue; i++) {
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
        if (dropdownStringList.first == "0") {
          selectedParamIndex--;
        }
        CreateData.selectedActionOrReactionData.parameterList[parameterList[index].name] = selectedParamIndex;
      }
      return CreateDropdown(
        index: index,
        onChanged: (String? newValue) {
          setState(() {
            selectedItems[index] = newValue ?? "";
            int selectedParamIndex = dropdownStringList.indexOf(selectedItems[index] ?? "") + 1;
            if (dropdownStringList.first == "0") {
              selectedParamIndex--;
            }
            CreateData.selectedActionOrReactionData.parameterList[parameterList[index].name] = selectedParamIndex;
            fillParameterFiledList();
          });
        },
        stringList: dropdownStringList,
        selectedItems: selectedItems,
        dropdownName: "${parameterList[index].name} : "
      );
    }
    if (parameterList[index].route.isEmpty) {
      textInputControler[index] = new TextEditingController(text: "");
      return Padding(
        padding: EdgeInsets.only(top: 20),
        child: Container(
          padding: EdgeInsets.all(0.0),
          constraints: BoxConstraints(maxHeight: 200.0, maxWidth: (MediaQuery.of(Globaldata.myContext).size.width - 40)),
          child: SingleChildScrollView(
            child: CustomClassicForm(
              hint: "${parameterList[index].name} : ",
              controller: textInputControler[index],
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

  @override
  Widget build(BuildContext context) {
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
              backgroundColor: Color.fromARGB(255, 34, 34, 34),
              text: AutoSizeText(
                "Continue",
                style: TextStyle(
                    fontSize: 30,
                    fontWeight: FontWeight.bold,
                    color: Colors.white),
              ),
              onPressed: () {
                for (int i = 0; i < textInputControler.length; i++) {
                  if (textInputControler[i] != null) {
                    String controlerText = textInputControler[i]!.text;
                    if (CreateData.selectedActionOrReactionData.parameter[i].type == "int") {
                      CreateData.selectedActionOrReactionData.parameterList[CreateData.selectedActionOrReactionData.parameter[i].name] = int.parse(controlerText);
                    } else {
                      CreateData.selectedActionOrReactionData.parameterList[CreateData.selectedActionOrReactionData.parameter[i].name] = controlerText;
                    }
                    developer.log("text input content : ${textInputControler[i]?.text}");
                  }
                }
                developer.log("parameter : ${CreateData.selectedActionOrReactionData.parameterList}");
                CreateData.selectOrUpdateWorkflow = true;
                Navigator.pop(Globaldata.myContext);
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

class CompleteActionOrReactionTop extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        SafeWebImage(
          url:
              "${Globaldata.domainName}${CreateData.serviceIconsList[Globaldata.firstService]}",
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
    return Column(
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
          value: selectedItems[index],
          hint: Text(
            selectedItems[index] ?? "",
            style: TextStyle(
              fontSize: 20,
              fontWeight: FontWeight.bold,
              color: Colors.white,
            )
          ),
          items: stringList.map((String item) {
            return DropdownMenuItem<String>(
              value: item,
              child: Text(
                item,
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
    );
  }
}
